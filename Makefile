build:
	go build -o bin/rwapigolang ./cmd

run: build
	./bin/rwapigolang

test:
	go test-v ./...