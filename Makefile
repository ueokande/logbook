GO_FILES=$(shell find -name '*.go' -not -name '*_test.go')
TARGET=logbook
TARGET_LINUX_AMD64=logbook_linux_amd64
TARGET_DARWIN_AMD64=logbook_darwin_amd64
TARGET_WINDOWS_AMD64=logbook_windows_amd64.exe

build: $(TARGET)
$(TARGET): $(GO_FILES)
	go build -o $(TARGET) .

binaries: $(TARGET_LINUX_AMD64) $(TARGET_DARWIN_AMD64) $(TARGET_WINDOWS_AMD64)
$(TARGET_LINUX_AMD64): $(GO_FILES)
	GOOS=linux GOARCH=amd64 go build -o $(TARGET_LINUX_AMD64) .
$(TARGET_DARWIN_AMD64): $(GO_FILES)
	GOOS=darwin GOARCH=amd64 go build -o $(TARGET_DARWIN_AMD64) .
$(TARGET_WINDOWS_AMD64): $(GO_FILES)
	GOOS=windows GOARCH=amd64 go build -o $(TARGET_WINDOWS_AMD64) .

test:
	go test -v -race ./...

lint:
	test -z "$$(gofmt -s -d . | tee /dev/stderr)"
	test -z "$$(golint ./...| tee /dev/stderr)"

clean:
	rm -rf $(TARGET) $(TARGET_LINUX_AMD64) $(TARGET_DARWIN_AMD64) $(TARGET_WINDOWS_AMD64)

.PHONY: build binaries test lint clean
