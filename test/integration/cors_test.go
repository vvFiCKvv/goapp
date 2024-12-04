package integration

import (
	integration "goapp/test/integration/base"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
)

func TestInvalidOrigin(t *testing.T) {
	teardownSuite := integration.Helper.SetupTest()
	defer teardownSuite()

	headers := http.Header{}
	testOrigin := "http://test.io"
	headers.Set("origin", testOrigin)
	_, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/goapp/ws", headers)
	if err == nil {
		t.Fatalf(`Origin %s should not be acceptable`, testOrigin)
	}
}

func TestValidOrigin(t *testing.T) {
	teardownSuite := integration.Helper.SetupTest()
	defer teardownSuite()
	headers := http.Header{}
	testOrigin := "http://mytestsite.io"
	headers.Set("origin", testOrigin)
	_, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/goapp/ws", headers)
	if err != nil {
		t.Fatalf(`Origin %s should be accepted, failed with error: %+v`, testOrigin, err)
	}
}

func TestSecondValidOrigin(t *testing.T) {
	teardownSuite := integration.Helper.SetupTest()
	defer teardownSuite()
	headers := http.Header{}
	testOrigin := "http://localhost:8080"
	headers.Set("origin", testOrigin)
	_, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/goapp/ws", headers)
	if err != nil {
		t.Fatalf(`Origin %s should be accepted, failed with error: %+v`, testOrigin, err)
	}
}
