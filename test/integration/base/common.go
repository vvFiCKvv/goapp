package integration

import (
	goapp "goapp/internal/app/server"
	"log"
	"os"
	"testing"
)

type IntegrationTest struct {
	testing.TB
}

func (tb *IntegrationTest) SetupTest() func() {
	monkeyPatchesToRestore := monkeyPatches()
	exitChannel := make(chan os.Signal, 1)
	serverOptions := goapp.ServerStartOptions{ExitChannel: exitChannel, UseProfiler: false}
	go func() {
		if err := goapp.Start(&serverOptions); err != nil {
			log.Fatalf("fatal: %+v\n", err)
		}
	}()
	return func() {
		for _, restoreCallback := range monkeyPatchesToRestore {
			restoreCallback()
		}
		exitChannel <- os.Kill
	}
}

var Helper = IntegrationTest{}
