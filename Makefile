clean:
	rm -rf ./bin

build:
	mkdir -p ./bin
	go build -o ./bin/issue-summoner main.go

test:
	go test -v ./...

coverage:
	go clean -testcache && go test -v -cover ./...


.PHONY: clean build run test coverage
