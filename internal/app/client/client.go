package client

import (
	"encoding/json"
	"goapp/internal/pkg/watcher"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var Port = 8080

var Success = func(connectionIndex int, iteration int, value string) {
	log.Printf(`[conn #%d] iteration: %d, value: %s`, connectionIndex, iteration, value)
}

var Fail = func(connectionIndex int, message string, err error) {
	log.Printf(`Error: #[conn #%d] %s, failed with error: %+v`, connectionIndex, message, err)
}

type ClientStartOptions struct {
	ParallelConnections int
	MessagesToSent      int
}

func Start(options *ClientStartOptions) error {

	w := sync.WaitGroup{}
	w.Add(options.ParallelConnections)
	for i := 0; i < options.ParallelConnections; i++ {
		go func(wg *sync.WaitGroup, index int) {
			defer wg.Done()
			connectAndReceiveData(index, options.MessagesToSent)
		}(&w, i)
	}
	w.Wait()
	log.Printf(`Finished`)
	return nil
}

func connectAndReceiveData(connectionIndex int, messagesToSent int) bool {
	connection, err := connect(connectionIndex)
	if err != nil {
		return false
	}
	defer connection.Close()

	messageChannel := make(chan string, messagesToSent)
	errorChannel := make(chan error, messagesToSent)

	bindReceiveChannels(connection, messageChannel, errorChannel)

	err = receiveMessages(messageChannel, errorChannel, connectionIndex, messagesToSent)
	if err != nil {
		return false
	}

	time.Sleep(time.Millisecond * 10)
	err = connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return false
	}

	time.Sleep(time.Millisecond * 10)
	return err == nil

}

func receiveMessages(messageChannel chan string, errorChannel chan error, connectionIndex int, messagesToSent int) error {
	receivedMessages := 0
	for {
		select {
		case message := <-messageChannel:

			var response watcher.Counter
			err := json.Unmarshal([]byte(message), &response)
			if err != nil {
				Fail(connectionIndex, `Failed to Unmarshal websocket response message`, err)
				return err
			}
			Success(connectionIndex, response.Iteration, response.Value)
			receivedMessages++
			if receivedMessages >= messagesToSent {
				return nil
			}

		case err := <-errorChannel:
			Fail(connectionIndex, `Failed to read a websocket message`, err)
			return err
		}
	}

}

func bindReceiveChannels(connection *websocket.Conn, messageChannel chan string, errorChannel chan error) {
	go func() {
		defer close(messageChannel)
		defer close(errorChannel)
		for {
			_, message, err := connection.ReadMessage()
			if err != nil {
				errorChannel <- err
			}
			messageChannel <- string(message[:])
		}
	}()
}

func connect(connectionIndex int) (*websocket.Conn, error) {
	headers := http.Header{}
	originUrl := "http://localhost:" + strconv.Itoa(Port)
	headers.Set("origin", originUrl)
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:"+strconv.Itoa(Port)+"/goapp/ws", headers)
	if err != nil {
		Fail(connectionIndex, `Connection`, err)
		return nil, err
	}
	return c, nil
}
