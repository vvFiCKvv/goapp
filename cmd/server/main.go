package main

import (
	"flag"
	goapp "goapp/internal/app/server"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

type MainArguments struct {
	useProfiler bool
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmsgprefix | log.Lshortfile)
}

func parseMainArguments() *MainArguments {
	profilePtr := flag.Bool("use-profiler", false, "Boolean to using pprof profiling")
	flag.Parse()

	return &MainArguments{
		useProfiler: *profilePtr,
	}
}

func main() {
	arguments := parseMainArguments()

	// Debug.
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	// Register signal handlers for exiting
	exitChannel := make(chan os.Signal, 1)
	signal.Notify(exitChannel, syscall.SIGINT, syscall.SIGTERM)

	// Start.
	serverOptions := goapp.ServerStartOptions{ExitChannel: exitChannel, UseProfiler: arguments.useProfiler}
	if err := goapp.Start(&serverOptions); err != nil {
		log.Fatalf("fatal: %+v\n", err)
	}
}
