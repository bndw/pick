PWD := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

GOPKG = github.com/bndw/pick
GOPATH = "$(CURDIR)/vendor:$(CURDIR)"

PICK_DIR = $(HOME)/.pick
BIN_DIR = /usr/local/bin
INSTALL = install

FOLDERS = $(shell find ./src -mindepth 1 -maxdepth 1 -type d)

all: build

install_hooks:
	@if test -d .git; then \
		for HOOK in githooks/???*; do \
			case $$HOOK in \
				*.sample|*~|*.swp) continue;; \
			esac; \
			if test -x $$HOOK; then \
				test ! -x .git/hooks/$${HOOK##*/} && echo "Installing git hook $${HOOK##*/}.."; \
				$(INSTALL) -c $$HOOK .git/hooks; \
			fi \
		done \
	fi

goget:
	GOPATH=$(GOPATH) go get github.com/rogpeppe/godeps
	GOPATH=$(GOPATH) $(CURDIR)/vendor/bin/godeps -u dependencies.tsv
	mkdir -p $(shell dirname "$(CURDIR)/vendor/src/$(GOPKG)")
	rm -f $(CURDIR)/vendor/src/$(GOPKG)
	ln -sf $(PWD)/src $(CURDIR)/vendor/src/$(GOPKG)

build: install_hooks goget
	GOPATH=$(GOPATH) go build -o bin/pick $(GOPKG)

test: goget
	GOPATH=$(GOPATH) go test -v $(FOLDERS)

install:
	@echo "Installing pick to $(BIN_DIR)/pick"
	$(INSTALL) -c bin/pick $(BIN_DIR)/pick

uninstall:
	rm -f $(BIN_DIR)/pick

fmt: gofmt

gofmt:
	GOPATH=$(GOPATH) go fmt $(FOLDERS)

config:
	@if [ ! -f "$(PICK_DIR)/config.toml" ]; then \
		OLD_UMASK=$(shell echo `umask`) ; \
		umask 077 ; \
		mkdir -p $(PICK_DIR) ; \
		$(INSTALL) -c -m 0600 config.toml.in $(PICK_DIR)/config.toml ; \
		umask $(OLD_UMASK) ; \
	fi

clean:
	rm -rf vendor/
	rm -rf bin/

.PHONY: all install_hooks goget build test install uninstall fmt gofmt config clean
