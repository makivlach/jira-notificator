BUILDPATH=$(CURDIR)
GO=$(shell which go)
GOINSTALL=$(GO) install
GOBINDATA=go-bindata
GOCLEAN=$(GO) clean
GOGET=$(GO) get

get:
	@$(GOGET) -u

build:
	@echo "start building..."
	$(GOBINDATA) -o ./cmd/jira-notificator-gui/bindata.go assets/
	$(GO) build ./cmd/jira-notificator-gui

all: get build