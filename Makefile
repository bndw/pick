PWD := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

GOPKG = github.com/bndw/pick
GOVENDOR = vendor
GOPATH = "$(CURDIR)/$(GOVENDOR)"

PICK_DIR = $(HOME)/.pick
BIN_DIR = /usr/local/bin
INSTALL = install

FOLDERS = $(shell find . -mindepth 1 -type d -not -path "*.git*" -not -path "./githooks*" -not -path "./$(GOVENDOR)*" -not -path "./Godeps*" -not -path "*bin*")

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

dependencies:
	@$(shell \
		cd $(GOVENDOR) ; \
		rm -rf src ; \
		find . -mindepth 3 -maxdepth 3 -path ./src -prune -o -type d -print | \
		sed -e 's/.\///' | \
		xargs -I{} sh -c ' \
			mkdir -p "src/`dirname {}`" ; \
			ln -sfn "../../../{}" "src/{}" ; \
		' \
	)
	@mkdir -p $(shell dirname $(GOVENDOR)/src/$(GOPKG))
	@ln -sfn ../../../.. $(GOVENDOR)/src/$(GOPKG)

build: install_hooks dependencies
	GOPATH=$(GOPATH) go build -o bin/pick .

test: dependencies
	GOPATH=$(GOPATH) go test -v $(FOLDERS)

install:
	@echo "Installing pick to $(BIN_DIR)/pick"
	$(INSTALL) -c bin/pick $(BIN_DIR)/pick

uninstall:
	rm -f $(BIN_DIR)/pick

fmt: gofmt

gofmt:
	GOPATH=$(GOPATH) go fmt $(FOLDERS)

govet:
	GOPATH=$(GOPATH) go tool vet $(FOLDERS)

config:
	@if [ ! -f "$(PICK_DIR)/config.toml" ]; then \
		OLD_UMASK=$(shell echo `umask`) ; \
		umask 077 ; \
		mkdir -p $(PICK_DIR) ; \
		$(INSTALL) -c -m 0600 config.toml.in $(PICK_DIR)/config.toml ; \
		umask $(OLD_UMASK) ; \
	fi

clean:
	rm -rf bin/

.PHONY: all install_hooks dependencies build test install uninstall fmt gofmt govet config clean
