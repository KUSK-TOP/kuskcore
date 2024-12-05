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

KUSKD_BINARY32 := kuskd-$(GOOS)_386
KUSKD_BINARY64 := kuskd-$(GOOS)_amd64

KUSKCLI_BINARY32 := kuskcli-$(GOOS)_386
KUSKCLI_BINARY64 := kuskcli-$(GOOS)_amd64

VERSION := $(shell awk -F= '/Version =/ {print $$2}' version/version.go | tr -d "\" ")

KUSKD_RELEASE32 := kuskd-$(VERSION)-$(GOOS)_386
KUSKD_RELEASE64 := kuskd-$(VERSION)-$(GOOS)_amd64

KUSKCLI_RELEASE32 := kuskcli-$(VERSION)-$(GOOS)_386
KUSKCLI_RELEASE64 := kuskcli-$(VERSION)-$(GOOS)_amd64

KUSK_RELEASE32 := kusk-$(VERSION)-$(GOOS)_386
KUSK_RELEASE64 := kusk-$(VERSION)-$(GOOS)_amd64

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
	@mkdir -p $@

binary: target/$(KUSKD_BINARY32) target/$(KUSKD_BINARY64) target/$(KUSKCLI_BINARY32) target/$(KUSKCLI_BINARY64)

ifeq ($(GOOS),windows)
release: binary
	cd target && cp -f $(KUSKD_BINARY32) $(KUSKD_BINARY32).exe
	cd target && cp -f $(KUSKCLI_BINARY32) $(KUSKCLI_BINARY32).exe
	cd target && md5sum $(KUSKD_BINARY32).exe $(KUSKCLI_BINARY32).exe > $(KUSK_RELEASE32).md5
	cd target && zip $(KUSK_RELEASE32).zip $(KUSKD_BINARY32).exe $(KUSKCLI_BINARY32).exe $(KUSK_RELEASE32).md5
	cd target && rm -f $(KUSKD_BINARY32) $(KUSKCLI_BINARY32) $(KUSKD_BINARY32).exe $(KUSKCLI_BINARY32).exe $(KUSK_RELEASE32).md5
	cd target && cp -f $(KUSKD_BINARY64) $(KUSKD_BINARY64).exe
	cd target && cp -f $(KUSKCLI_BINARY64) $(KUSKCLI_BINARY64).exe
	cd target && md5sum $(KUSKD_BINARY64).exe $(KUSKCLI_BINARY64).exe > $(KUSK_RELEASE64).md5
	cd target && zip $(KUSK_RELEASE64).zip $(KUSKD_BINARY64).exe $(KUSKCLI_BINARY64).exe $(KUSK_RELEASE64).md5
	cd target && rm -f $(KUSKD_BINARY64) $(KUSKCLI_BINARY64) $(KUSKD_BINARY64).exe $(KUSKCLI_BINARY64).exe $(KUSK_RELEASE64).md5
else
release: binary
	cd target && md5sum $(KUSKD_BINARY32) $(KUSKCLI_BINARY32) > $(KUSK_RELEASE32).md5
	cd target && tar -czf $(KUSK_RELEASE32).tgz $(KUSKD_BINARY32) $(KUSKCLI_BINARY32) $(KUSK_RELEASE32).md5
	cd target && rm -f $(KUSKD_BINARY32) $(KUSKCLI_BINARY32) $(KUSK_RELEASE32).md5
	cd target && md5sum $(KUSKD_BINARY64) $(KUSKCLI_BINARY64) > $(KUSK_RELEASE64).md5
	cd target && tar -czf $(KUSK_RELEASE64).tgz $(KUSKD_BINARY64) $(KUSKCLI_BINARY64) $(KUSK_RELEASE64).md5
	cd target && rm -f $(KUSKD_BINARY64) $(KUSKCLI_BINARY64) $(KUSK_RELEASE64).md5
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

target/$(KUSKD_BINARY32):
	CGO_ENABLED=0 GOARCH=386 go build $(BUILD_FLAGS) -o $@ cmd/kuskd/main.go

target/$(KUSKD_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/kuskd/main.go

target/$(KUSKCLI_BINARY32):
	CGO_ENABLED=0 GOARCH=386 go build $(BUILD_FLAGS) -o $@ cmd/kuskcli/main.go

target/$(KUSKCLI_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/kuskcli/main.go

test:
	@echo "====> Running go test"
	@go test $(PACKAGES)

benchmark:
	@go test -bench $(PACKAGES)

functional-tests:
	@go test -timeout=5m -tags="functional" ./test

ci: test

.PHONY: all kuskd kuskcli install target release release-all clean test benchmark functional-tests ci