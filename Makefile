.PHONY: lint
lint:
	go vet ./...
	staticcheck -checks="-SA1029" -f stylish ./...

.PHONY: build_win
build_win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o program.exe cmd/main.go

.PHONY: build_linux
build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o program cmd/main.go

.PHONY: swagger-docs
swagger-docs:
	swag init -g api.go -d ./internal/api,./internal/response,./internal/service -o ./docs