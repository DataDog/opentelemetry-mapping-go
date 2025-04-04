SHELL = /bin/bash
TOOLS_MOD_DIR := $(CURDIR)/internal/tools
GOTEST_OPT?= -race -timeout 2400s

.PHONY: install-tools
install-tools:
	cd $(TOOLS_MOD_DIR) && go install go.opentelemetry.io/build-tools/chloggen
	cd $(TOOLS_MOD_DIR) && go install go.opentelemetry.io/build-tools/multimod
	cd $(TOOLS_MOD_DIR) && go install github.com/frapposelli/wwhrd
	cd $(TOOLS_MOD_DIR) && go install github.com/golangci/golangci-lint/cmd/golangci-lint
	cd $(TOOLS_MOD_DIR) && go install golang.org/x/exp/cmd/apidiff@v0.0.0-20250106191152-7588d65b2ba8
	cd $(TOOLS_MOD_DIR) && go install gotest.tools/gotestsum
	cd $(TOOLS_MOD_DIR)/generate-license-file && go install .

FILENAME?=$(shell git branch --show-current).yaml
.PHONY: chlog-new
chlog-new:
	chloggen new --filename $(FILENAME)

.PHONY: chlog-validate
chlog-validate:
	chloggen validate

.PHONY: chlog-preview
chlog-preview:
	chloggen update --dry

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
	@$(MAKE) for-all CMD="go mod tidy -compat=1.23"

# Format code for all modules
.PHONY: fmt
fmt:
	@$(MAKE) for-all CMD="gofmt -w -s ./"

# Run unit test suite for all modules
.PHONY: test
test:
	@$(MAKE) for-all CMD="gotestsum -- $(GOTEST_OPT) ./..."

.PHONY: test-junit
test-junit:
	mkdir -p $(TOOLS_MOD_DIR)/junit
	COUNT=0
	set -e; for mod in $(GOMODULES); do \
  		COUNT=$$((COUNT+1)); \
		cd $$mod && gotestsum --junitfile $(CURDIR)/testresults/$$COUNT-junit.xml -- $(GOTEST_OPT) -coverprofile=coverage.cover -covermode=atomic ./... && cd -; \
	done

# Run linters for all modules
# Use 'make lint OPTS="--fix"' to autofix issues.
.PHONY: lint
lint:
	@$(MAKE) for-all CMD="golangci-lint run ./... $(OPTS)"

# Generate licenses file for compliance. 
.PHONY: gen-licenses
gen-licenses:
	generate-license-file

# Do PR for preparing a release
.PHONY: prerelease
prerelease:
	multimod verify && multimod prerelease -m pkgs

# Push tags 
.PHONY: push-tags
push-tags:
	multimod verify
	set -e; for tag in `multimod tag -m pkgs -c HEAD --print-tags | grep -v "Using" `; do \
		echo "pushing tag $${tag}"; \
		git push git@github.com:DataDog/opentelemetry-mapping-go.git $${tag}; \
	done;

APIHEADERS := internal/apidiff-data

.PHONY: apidiff-generate
apidiff-generate:
	set -e; for mod in $(GOMODULES); do \
		./internal/scripts/apidiff-generate.sh $$mod $(APIHEADERS); \
	done

.PHONY: apidiff-compare
apidiff-compare:
	set -e; for mod in $(GOMODULES); do \
		./internal/scripts/apidiff-compare.sh $$mod $(APIHEADERS); \
	done
