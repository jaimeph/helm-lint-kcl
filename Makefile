lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.62.2 golangci-lint run -v
