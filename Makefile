BINARY = go-anonymize-mysqldump
BUILDDIR = ./build

all: clean vet fmt lint test build

build:
	gox -os="linux" -os="darwin" -os="windows" -arch="amd64" -arch="386" -output="${BUILDDIR}/${BINARY}_{{.OS}}_{{.Arch}}"
	gzip build/*

vet:
	go get -v -t -d ./...

fmt:
	gofmt -s -l . | grep -v vendor | tee /dev/stderr

lint:
	golint ./... | grep -v vendor | tee /dev/stderr

test:
	LOG_LEVEL=debug go test -v ./...
	go test -bench .

test-watch:
	fswatch -0 *.go | xargs -0 -L 1 -I % sh -c 'LOG_LEVEL=debug go test -v ./...'

run-test-watch:
	fswatch -0 *.go | xargs -0 -L 1 -I % sh -c 'LOG_LEVEL=debug go test -v -run $(TEST)'

clean:
	rm -rf ${BUILDDIR}

.PHONY: all clean vet fmt lint test build
