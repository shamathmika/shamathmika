.PHONY: render test

render:
	go run ./cmd/pet render

test:
	go test ./...
