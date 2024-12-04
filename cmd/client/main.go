package main

import (
	"flag"
	client "goapp/internal/app/client"
	"log"
	"os"
)

type MainArguments struct {
	n int
}

func parseMainArguments() *MainArguments {
	fs := flag.NewFlagSet("client", flag.ContinueOnError)
	n := fs.Int("n", 1, "Int64 representing the number of multiple parallel connections")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Println(err) // prints an error message
	}

	return &MainArguments{
		n: *n,
	}
}

func main() {
	arguments := parseMainArguments()

	clientOptions := client.ClientStartOptions{ParallelConnections: arguments.n, MessagesToSent: 1}

	if err := client.Start(&clientOptions); err != nil {
		log.Fatalf("fatal: %+v\n", err)
	}
}
