GO ?= go
APP ?= vtrix
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
PREFIX ?= /usr/local

VTRIX_BASE_URL ?= https://vtrix.ai
VTRIX_MODELS_URL ?= https://seacloud-cloud-model-spec.api.seaart.ai
VTRIX_GENERATION_URL ?= $(VTRIX_BASE_URL)
VTRIX_SKILLHUB_URL ?= https://skill-hub.vtrix.ai/api/v1

LDFLAGS := -s -w \
	-X github.com/VtrixAI/vtrix-cli/internal/buildinfo.Version=$(VERSION) \
	-X github.com/VtrixAI/vtrix-cli/internal/auth.BaseURL=$(VTRIX_BASE_URL) \
	-X github.com/VtrixAI/vtrix-cli/internal/models.BaseURL=$(VTRIX_MODELS_URL) \
	-X github.com/VtrixAI/vtrix-cli/internal/generation.BaseURL=$(VTRIX_GENERATION_URL) \
	-X github.com/VtrixAI/vtrix-cli/internal/skillhub.BaseURL=$(VTRIX_SKILLHUB_URL)

.PHONY: build install uninstall clean

build:
	$(GO) build -ldflags "$(LDFLAGS)" -o $(APP) .

install: build
	install -d $(PREFIX)/bin
	install -m755 $(APP) $(PREFIX)/bin/$(APP)
	@echo "installed $(APP) to $(PREFIX)/bin/$(APP)"

uninstall:
	rm -f $(PREFIX)/bin/$(APP)

clean:
	rm -f "$(APP)"
