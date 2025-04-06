VERSION:=$(shell git log --date=short --pretty=format:'%cd-%h' -n 1)

.PHONY: start
start:
	@echo "Running SnipAndNeat application"
	@echo "Configuration file: .env"
	@echo "To change configuration, please edit .env file"
	@echo "For more information, please visit https://github.com/ilius1995/SnipAndNeat"
	@echo ""

	@env -f .env go run -ldflags "-X main.build=$(shell git rev-parse HEAD)" main.go

test:
	@./go.test.sh
.PHONY: test

coverage:
	@./go.coverage.sh
.PHONY: coverage

generate:
	go generate ./...
.PHONY: generate

check_generated: generate
	git diff --exit-code
.PHONY: check_generated

build:
	GOPRIVATE=$(GOPRIVATE) CGO_ENABLED=1 go build \
		-v \
		-ldflags "-w -s -extldflags '-static' -X main.VERSION=$(VERSION)" \
		-tags 'netgo std static_all' \
		-o output ./cmd/SnipAndNeat/main.go