get:
	go get -t ./...
.PHONY: get

install:
	go install github.com/Sean-Clarke/go-snake-go
.PHONY: install

run: install
	go-snake-go server
.PHONY: run

test:
	go test ./...
.PHONY: test

fmt:
	gofmt -l -s -w .