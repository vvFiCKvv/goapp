.DEFAULT_GOAL := all

.PHONY: all
all: clean goapp client

.PHONY: goapp
goapp:
	mkdir -p bin
	go build -o bin ./cmd/server

.PHONY: client
client:
	mkdir -p bin
	go build -o bin ./cmd/client

.PHONY: stress_sockets
stress_sockets:
	./bin/server -use-profiler&
	k6 run ./test/stress/sockets.k6.ts
	sleep 4
	kill $$(ps aux | grep './bin/server' | head -n 1| awk '{print $$2}')
	sleep 4
	go tool pprof --web ./.pprof/mem.pprof

.PHONY: clean
clean:
	go clean
	rm -f bin/*

.PHONY: test
test:
	go test ./...
