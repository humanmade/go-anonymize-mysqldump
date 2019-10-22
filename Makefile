BINARY = go-anonymize-mysqldump
BUILDDIR = ./build

all: clean vet fmt lint test build

build:
	gox -os="linux" -os="darwin" -arch="amd64" -output="${BUILDDIR}/${BINARY}_{{.OS}}_{{.Arch}}"
	gzip build/*

vet:
	go get -v -t -d ./...

fmt:
	gofmt -s -l . | grep -v vendor | tee /dev/stderr

lint:
	golint ./... | grep -v vendor | tee /dev/stderr

test:
	go test -v ./...
	go test -bench .

test-watch:
	fswatch -0 *.go | xargs -0 -L 1 -I % go test -v ./...

clean:
	rm -rf ${BUILDDIR}

.PHONY: all clean vet fmt lint test build
