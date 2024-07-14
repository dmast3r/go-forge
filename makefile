build:
	mkdir -p bin
	go build -o bin/app app/main/server.go

run:
	./bin/app

.PHONY: build run