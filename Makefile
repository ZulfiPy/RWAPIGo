build:
	go build -o bin/rwapigolang

run: build
	./bin/rwapigolang

test:
	go test-v ./...