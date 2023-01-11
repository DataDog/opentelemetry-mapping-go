TOOLS_MOD_DIR := ./internal/tools

.PHONY: install-tools
install-tools:
	cd $(TOOLS_MOD_DIR) && go install go.opentelemetry.io/build-tools/chloggen


FILENAME?=$(shell git branch --show-current).yaml
.PHONY: chlog-new
chlog-new:
	chloggen new --filename $(FILENAME)

GOMODULES := $(shell find . -type f -name "go.mod" -exec dirname {} \; | sort | egrep  '^./' )

.PHONY: $(GOMODULES)
$(GOMODULES):
	@echo "Running '$(CMD)' in module '$@'"
	cd $@ && $(CMD)

# Run CMD for all modules
.PHONY: for-all
for-all: $(GOMODULES)

# Tidy go.mod/go.sum for all modules
.PHONY: tidy
tidy:
	@$(MAKE) for-all CMD="go mod tidy -compat=1.18"

# Format code for all modules
.PHONY: fmt
fmt:
	@$(MAKE) for-all CMD="gofmt -w -s ./"

# Run unit test suite for all modules
.PHONY: test
test:
	@$(MAKE) for-all CMD="go test -race -timeout 300s ./..."
