clean:
	rm -rf ./bin

build:
	mkdir -p ./bin
	go build -o ./bin/issue-summoner main.go

.PHONY: clean build run