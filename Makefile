PWD := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

GOPKG = github.com/bndw/pick
GOPATH = "$(CURDIR)/vendor:$(CURDIR)"

PICK_DIR = $(HOME)/.pick
BIN_DIR = /usr/local/bin
INSTALL = install

FOLDERS = $(shell find . -mindepth 1 -maxdepth 1 -type d -not -path "*.git" -not -path "*vendor" -not -path "*bin")

all: build

goget:
	GOPATH=$(GOPATH) go get github.com/rogpeppe/godeps
	GOPATH=$(GOPATH) $(CURDIR)/vendor/bin/godeps -u dependencies.tsv
	mkdir -p $(shell dirname "$(CURDIR)/vendor/src/$(GOPKG)")
	rm -f $(CURDIR)/vendor/src/$(GOPKG)
	ln -sf $(PWD) $(CURDIR)/vendor/src/$(GOPKG)

build: goget
	GOPATH=$(GOPATH) go build -o bin/pick .

test: goget
	GOPATH=$(GOPATH) go test -v $(FOLDERS)

install:
	@echo "Installing pick to $(BIN_DIR)/pick"
	$(INSTALL) -c bin/pick $(BIN_DIR)/pick

uninstall:
	rm -f $(BIN_DIR)/pick

config:
	@if [ ! -f "$(PICK_DIR)/config.toml" ]; then \
		$(INSTALL) -c -m 0600 config.toml.in $(PICK_DIR)/config.toml ; \
	fi

clean:
	rm -rf vendor/
	rm -rf bin/

.PHONY: all goget build test install uninstall config clean
