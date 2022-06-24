.DEFAULT_GOAL := run

run:
	@go mod tidy -v
	@go run cmd/url-shrtnr/main.go
