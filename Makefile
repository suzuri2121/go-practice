-include .env
export

.PHONY: test run

test:
	go test -v -count=1 ./...

run:
	go run ./cmd/api