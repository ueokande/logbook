GO_FILES=$(shell find -name '*.go' -not -name '*_test.go')
TARGET=logbook

build: $(TARGET)
$(TARGET): $(GO_FILES)
	go build .

test:
	go test -v -race ./...

lint:
	test -z "$$(gofmt -s -d . | tee /dev/stderr)"
	test -z "$$(golint ./...| tee /dev/stderr)"

clean:
	rm -rf $(TARGET)

.PHONY: build test lint clean
