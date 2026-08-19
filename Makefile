.PHONY: wire-build generate-docs run test

wire-build:
	@echo "Running wire in the di folder..."
	@wiregenx --root ./src --out ../di/wire/provider.go
	@cd di/wire/ && wire && mv wire_gen.go ../../cmd/container/container.go

generate-docs:
	@echo "Generating documentation..."
	@swag init -g ./cmd/http/main.go -o ./src/entrypoint/presentation/docs
	@echo "Documentation generated successfully."

run:
	@echo "Running the application..."
	@go run ./cmd/http/main.go

test:
	@echo "Running tests..."
	@go test ./... -count=1 -cover -coverpkg=./src/...
