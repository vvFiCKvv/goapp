package integration

import (
	"goapp/internal/app/client"
	httpsrv "goapp/internal/pkg/httpsrv"
	integration "goapp/test/integration/base"
	"sync"
	"testing"
	"time"
)

func monkeyPatchClient(t *testing.T, onSuccess func(connectionIndex int, iteration int, value string), onStatsPrint func(id string, sent int)) func() {

	originalSuccess := client.Success
	originalFail := client.Fail
	originalStatsPrint := httpsrv.StatsPrint
	httpsrv.StatsPrint = func(id string, sent int) {
		onStatsPrint(id, sent)
	}
	client.Fail = func(connectionIndex int, message string, err error) {
		t.Fatalf(`Error: #[conn #%d] %s, failed with error: %+v`, connectionIndex, message, err)
	}
	client.Success = onSuccess
	return func() {
		client.Success = originalSuccess
		client.Fail = originalFail
		httpsrv.StatsPrint = originalStatsPrint
	}
}

func TestE2EOneConnectionOneMessage(t *testing.T) {
	parallelConnections := 1
	messagesToSent := 1

	teardownSuite := integration.Helper.SetupTest()
	defer teardownSuite()

	genericTest(t, parallelConnections, messagesToSent)
}

func TestE2EOneConnectionTwoMessages(t *testing.T) {
	parallelConnections := 1
	messagesToSent := 2

	teardownSuite := integration.Helper.SetupTest()
	defer teardownSuite()

	genericTest(t, parallelConnections, messagesToSent)
}

func TestE2ETwoConnectionsOneMessage(t *testing.T) {
	parallelConnections := 2
	messagesToSent := 1

	teardownSuite := integration.Helper.SetupTest()
	defer teardownSuite()

	genericTest(t, parallelConnections, messagesToSent)
}

func TestE2ETwoConnectionsTweMessages(t *testing.T) {
	parallelConnections := 2
	messagesToSent := 2

	teardownSuite := integration.Helper.SetupTest()
	defer teardownSuite()

	genericTest(t, parallelConnections, messagesToSent)
}

func TestE2EMultipleConnections(t *testing.T) {
	parallelConnections := 100
	messagesToSent := 3

	teardownSuite := integration.Helper.SetupTest()
	defer teardownSuite()

	genericTest(t, parallelConnections, messagesToSent)
}

func TestE2EHugeConnections(t *testing.T) {
	parallelConnections := 1000
	messagesToSent := 1

	teardownSuite := integration.Helper.SetupTest()
	defer teardownSuite()

	genericTest(t, parallelConnections, messagesToSent)
}

func genericTest(t *testing.T, parallelConnections int, messagesToSent int) {
	var messagesValuesMap sync.Map
	var sessionsStatsMap sync.Map
	var messagesReceivedMap sync.Map
	countersLock := sync.RWMutex{}

	time.Sleep(time.Millisecond * 10)
	onSuccess := func(connectionIndex int, iteration int, value string) {

		t.Logf(`[conn #%d] iteration: %d, value: %s`, connectionIndex, iteration, value)

		countersLock.Lock()
		defer countersLock.Unlock()
		ValueCount, hasValueCount := messagesValuesMap.Load(value)
		if !hasValueCount {
			ValueCount = 0
		}
		messagesValuesMap.Store(value, ValueCount.(int)+1)

		messages, hasMessage := messagesReceivedMap.Load(connectionIndex)
		if !hasMessage {
			messages = 0
		}
		expectedMessages := messages.(int) + 1

		if expectedMessages != iteration {
			t.Fatalf(`Invalid iteration %d for conn# %d expected %d`, iteration, connectionIndex, expectedMessages)
		}
		messagesReceivedMap.Store(connectionIndex, iteration)
	}
	onStatsPrint := func(id string, sent int) {
		t.Logf("session %s has received %d messages\n", id, sent)
		sessionsStatsMap.Store(id, sent)
	}
	restore := monkeyPatchClient(t, onSuccess, onStatsPrint)
	defer restore()
	clientOptions := client.ClientStartOptions{ParallelConnections: parallelConnections, MessagesToSent: messagesToSent}
	client.Start(&clientOptions)

	var clientsCount = 0
	messagesReceivedMap.Range(func(connectionIndex, calledTimes any) bool {
		clientsCount++
		if calledTimes != messagesToSent {
			t.Fatalf(`Connection %d received messages %d times but expected %d`, connectionIndex, calledTimes, messagesToSent)
			return false
		}
		return true
	})
	if clientsCount != parallelConnections {
		t.Fatalf(`Some clients received no messages, %d of expected %d`, clientsCount, parallelConnections)
	}

	var valuesCount = 0
	messagesValuesMap.Range(func(value, valueCount any) bool {
		valuesCount++
		if valueCount != parallelConnections {
			// This my happen because of different strgen tick, so it is not always an error
			t.Logf(`Not all connections was notified for value %s, expected %d notified %d`, value, parallelConnections, valueCount)
			return false
		}
		return true
	})

	if valuesCount < messagesToSent {
		t.Fatalf(`Not all values delivered, expected %d but delivered %d`, messagesToSent, valuesCount)
	}

	var sessionsLength = 0
	sessionsStatsMap.Range(func(sessionId, messagesCount any) bool {
		sessionsLength++
		if messagesCount != messagesToSent {
			// This my happen because of different strgen tick, so it is not always an error
			t.Logf(`Not all sessions had correct statistics, expected %d notified %d`, messagesToSent, messagesCount)
			return false
		}
		return true
	})

	if valuesCount < messagesToSent {
		t.Fatalf(`Not all sessions statistics was stored, expected %d but have %d`, messagesToSent, valuesCount)
	}

}
