.PHONY: render preview test

render:
	go run ./cmd/pet render

preview:
	go run ./cmd/pet preview

test:
	go test ./...
