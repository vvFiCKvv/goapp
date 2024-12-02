.DEFAULT_GOAL := goapp

.PHONY: all
all: clean goapp

.PHONY: goapp
goapp:
	mkdir -p bin
	go build -o bin ./...

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
