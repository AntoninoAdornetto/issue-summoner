clean:
	rm -rf ./bin

build:
	mkdir -p ./bin
	go build -o ./bin/issue-summoner main.go

test:
	go test -v ./...

coverage:
	go clean -testcache && go test -coverprofile=coverage/coverage.out ./... && go tool cover -html=coverage/coverage.out -o=coverage/coverage.html


.PHONY: clean build run test coverage
