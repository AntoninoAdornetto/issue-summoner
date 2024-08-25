all: clean build 

clean:
	@rm -rf ./bin

build:
	@mkdir -p ./bin
	go build -o ./bin/issue-summoner main.go

test:
	@go test -v ./...

bench:
	@go test -v ./... -bench=. -benchtime=10s -run=^#

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

coverage:
	@go clean -testcache && go test -coverprofile=coverage/coverage.out ./... && go tool cover -html=coverage/coverage.out -o=coverage/coverage.html


.PHONY: all clean build test coverage lint
