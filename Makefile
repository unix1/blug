.PHONY: all test build install clean

all: test

test:
	go test ./...

build:
	go build -o blug .

install:
	go install .

clean:
	rm -f blug
