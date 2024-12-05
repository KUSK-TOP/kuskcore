ifndef GOOS
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	GOOS := darwin
else ifeq ($(UNAME_S),Linux)
	GOOS := linux
else
	GOOS := windows
endif
endif

PACKAGES    := $(shell go list ./... | grep -v '/lib/')

BUILD_FLAGS := -ldflags "-X github.com/KUSK-TOP/kuskcore/version.GitCommit=`git rev-parse HEAD`"

BYTOMD_BINARY32 := kuskd-$(GOOS)_386
BYTOMD_BINARY64 := kuskd-$(GOOS)_amd64

BYTOMCLI_BINARY32 := kuskcli-$(GOOS)_386
BYTOMCLI_BINARY64 := kuskcli-$(GOOS)_amd64

VERSION := $(shell awk -F= '/Version =/ {print $$2}' version/version.go | tr -d "\" ")

BYTOMD_RELEASE32 := kuskd-$(VERSION)-$(GOOS)_386
BYTOMD_RELEASE64 := kuskd-$(VERSION)-$(GOOS)_amd64

BYTOMCLI_RELEASE32 := kuskcli-$(VERSION)-$(GOOS)_386
BYTOMCLI_RELEASE64 := kuskcli-$(VERSION)-$(GOOS)_amd64

BYTOM_RELEASE32 := kusk-$(VERSION)-$(GOOS)_386
BYTOM_RELEASE64 := kusk-$(VERSION)-$(GOOS)_amd64

all: test target release-all install

kuskd:
	@echo "Building kuskd to cmd/kuskd/kuskd"
	@go build $(BUILD_FLAGS) -o cmd/kuskd/kuskd cmd/kuskd/main.go

kuskcli:
	@echo "Building kuskcli to cmd/kuskcli/kuskcli"
	@go build $(BUILD_FLAGS) -o cmd/kuskcli/kuskcli cmd/kuskcli/main.go

install:
	@echo "Installing kuskd and kuskcli to $(GOPATH)/bin"
	@go install ./cmd/kuskd
	@go install ./cmd/kuskcli

target:
	mkdir -p $@

binary: target/$(BYTOMD_BINARY32) target/$(BYTOMD_BINARY64) target/$(BYTOMCLI_BINARY32) target/$(BYTOMCLI_BINARY64)

ifeq ($(GOOS),windows)
release: binary
	cd target && cp -f $(BYTOMD_BINARY32) $(BYTOMD_BINARY32).exe
	cd target && cp -f $(BYTOMCLI_BINARY32) $(BYTOMCLI_BINARY32).exe
	cd target && md5sum  $(BYTOMD_BINARY32).exe $(BYTOMCLI_BINARY32).exe >$(BYTOM_RELEASE32).md5
	cd target && zip $(BYTOM_RELEASE32).zip  $(BYTOMD_BINARY32).exe $(BYTOMCLI_BINARY32).exe $(BYTOM_RELEASE32).md5
	cd target && rm -f  $(BYTOMD_BINARY32) $(BYTOMCLI_BINARY32)  $(BYTOMD_BINARY32).exe $(BYTOMCLI_BINARY32).exe $(BYTOM_RELEASE32).md5
	cd target && cp -f $(BYTOMD_BINARY64) $(BYTOMD_BINARY64).exe
	cd target && cp -f $(BYTOMCLI_BINARY64) $(BYTOMCLI_BINARY64).exe
	cd target && md5sum  $(BYTOMD_BINARY64).exe $(BYTOMCLI_BINARY64).exe >$(BYTOM_RELEASE64).md5
	cd target && zip $(BYTOM_RELEASE64).zip  $(BYTOMD_BINARY64).exe $(BYTOMCLI_BINARY64).exe $(BYTOM_RELEASE64).md5
	cd target && rm -f  $(BYTOMD_BINARY64) $(BYTOMCLI_BINARY64)  $(BYTOMD_BINARY64).exe $(BYTOMCLI_BINARY64).exe $(BYTOM_RELEASE64).md5
else
release: binary
	cd target && md5sum  $(BYTOMD_BINARY32) $(BYTOMCLI_BINARY32) >$(BYTOM_RELEASE32).md5
	cd target && tar -czf $(BYTOM_RELEASE32).tgz  $(BYTOMD_BINARY32) $(BYTOMCLI_BINARY32) $(BYTOM_RELEASE32).md5
	cd target && rm -f  $(BYTOMD_BINARY32) $(BYTOMCLI_BINARY32) $(BYTOM_RELEASE32).md5
	cd target && md5sum  $(BYTOMD_BINARY64) $(BYTOMCLI_BINARY64) >$(BYTOM_RELEASE64).md5
	cd target && tar -czf $(BYTOM_RELEASE64).tgz  $(BYTOMD_BINARY64) $(BYTOMCLI_BINARY64) $(BYTOM_RELEASE64).md5
	cd target && rm -f  $(BYTOMD_BINARY64) $(BYTOMCLI_BINARY64) $(BYTOM_RELEASE64).md5
endif

release-all: clean
	GOOS=darwin  make release
	GOOS=linux   make release
	GOOS=windows make release

clean:
	@echo "Cleaning binaries built..."
	@rm -rf cmd/kuskd/kuskd
	@rm -rf cmd/kuskcli/kuskcli
	@rm -rf target
	@rm -rf $(GOPATH)/bin/kuskd
	@rm -rf $(GOPATH)/bin/kuskcli
	@echo "Cleaning temp test data..."
	@rm -rf test/pseudo_hsm*
	@rm -rf blockchain/pseudohsm/testdata/pseudo/
	@echo "Cleaning sm2 pem files..."
	@rm -rf crypto/sm2/*.pem
	@echo "Done."

target/$(BYTOMD_BINARY32):
	CGO_ENABLED=0 GOARCH=386 go build $(BUILD_FLAGS) -o $@ cmd/kuskd/main.go

target/$(BYTOMD_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/kuskd/main.go

target/$(BYTOMCLI_BINARY32):
	CGO_ENABLED=0 GOARCH=386 go build $(BUILD_FLAGS) -o $@ cmd/kuskcli/main.go

target/$(BYTOMCLI_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/kuskcli/main.go

test:
	@echo "====> Running go test"
	@go test $(PACKAGES)

benchmark:
	@go test -bench $(PACKAGES)

functional-tests:
	@go test -timeout=5m -tags="functional" ./test 

ci: test

.PHONY: all target release-all clean test benchmark
