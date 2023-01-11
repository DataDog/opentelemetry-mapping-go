TOOLS_MOD_DIR := ./internal/tools

.PHONY: install-tools
install-tools:
	cd $(TOOLS_MOD_DIR) && go install go.opentelemetry.io/build-tools/chloggen


FILENAME?=$(shell git branch --show-current).yaml
.PHONY: chlog-new
chlog-new:
	chloggen new --filename $(FILENAME)
