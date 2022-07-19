.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -race -timeout 30s ./ ...

.DEFAULT_GOAL := build

# пока хз с MAKEFILE(на винде не робит) -  go build -v ./cmd/apiserver
