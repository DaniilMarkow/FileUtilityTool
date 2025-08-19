.PHONY: all test clean

INPUT_FILE ?= testdata/expected/filtered.csv
OUTPUT_FILE ?= testdata/expected/result.json

all: build
	./fileutil -input=$(INPUT_FILE) -output=$(OUTPUT_FILE) -sort=population

build:
	go build -o fileutil .

test:
	mkdir -p coverage
	go clean -testcache
	go test -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/report.html

clean:
	rm -f fileutil
	rm -f $(OUTPUT_FILE)
	rm -f coverage/coverage.out coverage/report.html
	rm -rf coverage