.PHONY: all build run test clean

build:
	go build -o mockfactory.exe ./cmd/mockfactory

run: build
	./mockfactory.exe

test:
	go test -v ./...

clean:
	rm -f mockfactory.exe

all: build