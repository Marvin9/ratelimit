.PHONY: dev test

dev:
	clear && go run main.go

test:
	clear && go test ./...