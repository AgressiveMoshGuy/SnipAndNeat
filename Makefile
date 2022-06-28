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