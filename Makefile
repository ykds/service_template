.PHONY: build
build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o program cmd/main.go

.PHONY: swagger-docs
swagger-docs:
	swag init -g api.go -d ./internal/api,./internal/response,./internal/service -o ./docs