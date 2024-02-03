all: dev

install:
	rm -f go.mod
	go mod init github.com/h-abranches-dev/daily-expenses-be

build:
	go build -o service

dev: build
	./service

clean:
	rm -f service
