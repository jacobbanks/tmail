.PHONY: build test clean install lint

build:
	go build -o tmail main.go

install:
	go install

test:
	go test ./...

lint:
	go vet ./...
	go fmt ./...

clean:
	rm -f tmail

run:
	go run main.go