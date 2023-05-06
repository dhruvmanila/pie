# By default, run the 'build' target.
.DEFAULT_GOAL = build

BIN := pie

.PHONY: build
build: ## Build the binary
	go build -v -o ./$(BIN) .

.PHONY: clean
clean: ## Remove build artifacts
	rm -f ./$(BIN)

.PHONY: tidy
tidy: ## Tidy the module
	go mod tidy

.PHONY: test
test: ## Run the tests
	go test ./... -v

HELP_FORMAT = "  \033[36m%-10s\033[0m %s\n"
.PHONY: help
help: ## Display this help message
	@# Print everything matching "target: ## magic comments"
	@echo "Valid targets:"
	@grep -E "^[^ ]+:.* ## .*" $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*? ## "}; \
				{printf $(HELP_FORMAT), $$1, $$2}'
