.PHONY: test run

test:
	go test ./...

run:
	go run ./cmd/contractlab ./examples/api.json
