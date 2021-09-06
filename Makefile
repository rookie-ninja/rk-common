
.PHONY: all
all: buf test lint readme fmt

.PHONY: buf
buf:
	@echo "running buf..."
	@buf generate

.PHONY: lint
lint:
	@echo "running golangci-lint..."
	@golangci-lint run 2>&1

.PHONY: test
test:
	@echo "running go test..."
	@go test -race ./... 2>&1

.PHONY: fmt
fmt:
	@echo "format go project..."
	@gofmt -s -w . 2>&1

.PHONY: readme
readme:
	@echo "running doctoc..."
	@doctoc . 2>&1