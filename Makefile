.PHONY: vet build clean install

# Prepend our _vendor directory to the system GOPATH
# so that import path resolution will prioritize
# our third party snapshots.
GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

default: build

vet:
	go vet ./cmd/main.go

build: vet
	go build -v -o ./bin/main ./cmd/main.go

clean:
	rm -rf ./vendor

install:
	glide up

