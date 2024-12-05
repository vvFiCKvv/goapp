package goapp

import (
	"fmt"
	"goapp/internal/pkg/httpsrv"
	"goapp/internal/pkg/strgen"
	"log"
	"os"
	"runtime"

	"github.com/pkg/profile"
)

type ServerStartOptions struct {
	ExitChannel chan os.Signal
	UseProfiler bool
}

func Start(options *ServerStartOptions) error {
	var (
		strChan = make(chan string, 100) // String channel with max parallel counter processes.
		strCli  = strgen.New(strChan)    // String generator.
		httpSrv = httpsrv.New(strChan)   // HTTP server.
	)

	// Start String Generator.
	if err := strCli.Start(); err != nil {
		return fmt.Errorf("failed to start string generator: %w", err)
	}
	defer strCli.Stop()

	// Start HTTP server.
	if err := httpSrv.Start(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	defer httpSrv.Stop()

	log.Println("GoApp Started")
	defer log.Println("GoApp Stopped")

	<-options.ExitChannel

	if options.UseProfiler {
		// we need to write the memory profile before the server defers, to see what remains in the heap
		writeMemoryProfile("./.pprof/")
	}

	return nil
}

func writeMemoryProfile(path string) {
	runtime.GC()
	profile.Start(profile.MemProfile, profile.ProfilePath(path)).Stop()
}
