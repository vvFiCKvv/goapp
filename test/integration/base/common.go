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

	exitChannel := make(chan os.Signal, 1)
	serverOptions := goapp.ServerStartOptions{ExitChannel: exitChannel, UseProfiler: false}
	go func() {
		if err := goapp.Start(&serverOptions); err != nil {
			log.Fatalf("fatal: %+v\n", err)
		}
	}()
	monkeyPatchesToRestore := monkeyPatches()
	return func() {
		for _, restoreCallback := range monkeyPatchesToRestore {
			restoreCallback()
		}
		exitChannel <- os.Kill
	}
}
