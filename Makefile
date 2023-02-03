OUTPUT ?= fiber-server

.DEFAULT_GOAL := build

build:
	CGO_ENABLED=0 go build -o $(OUTPUT)

all: build
