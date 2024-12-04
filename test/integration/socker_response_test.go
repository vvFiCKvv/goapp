package integration

import (
	"encoding/json"
	"goapp/internal/pkg/watcher"
	integration "goapp/test/integration/base"
	"net/http"
	"regexp"
	"testing"

	"github.com/gorilla/websocket"
)

func TestValidWSResponse(t *testing.T) {
	teardownSuite := integration.Helper.SetupTest()
	defer teardownSuite()
	headers := http.Header{}
	testOrigin := "http://mytestsite.io"
	headers.Set("origin", testOrigin)
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/goapp/ws", headers)
	if err != nil {
		t.Fatalf(`Connection failed with error: %+v`, err)
	}
	defer c.Close()

	messageChannel := make(chan string)
	errorChannel := make(chan error)
	go func() {
		defer close(messageChannel)
		defer close(errorChannel)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				errorChannel <- err
				return
			}
			messageChannel <- string(message[:])
		}
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte(""))
	if err != nil {
		t.Fatalf(`Failed to send a websocket message: %+v`, err)
		return
	}
	for {
		select {
		case message := <-messageChannel:
			t.Logf(`Response Message: %s`, message)

			var response watcher.Counter
			err := json.Unmarshal([]byte(message), &response)
			if err != nil {
				t.Fatalf(`Failed to Unmarshal websocket response message: %+v`, err)
			}
			if response.Iteration != 1 {
				t.Fatalf(`websocket responded with wrong iteration: %d`, response.Iteration)
			}

			r, _ := regexp.Compile("^[0-9A-F]+$")

			if r.MatchString(response.Value) == false {
				t.Fatalf(`websocket responded with wrong value: %s`, response.Value)
			}

			return
		case err = <-errorChannel:
			t.Fatalf(`Failed to read a websocket message: %+v`, err)
			return
		}
	}
}
