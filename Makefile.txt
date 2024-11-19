ifndef GOOS
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	GOOS := darwin
else ifeq ($(UNAME_S),Linux)
	GOOS := linux
else
$(error "$$GOOS is not defined. If you are using Windows, try to re-make using 'GOOS=windows make ...' ")
endif
endif

PACKAGES    := $(shell go list ./... | grep -v '/vendor/' )
PACKAGES += 'core/mining/tensority/go_algorithm'

BUILD_FLAGS := -ldflags "-X core/version.GitCommit=`git rev-parse HEAD`"

#KUSKD_BINARY32 := kuskd-$(GOOS)_386
KUSKD_BINARY64 := kuskd-$(GOOS)_amd64

#KUSKCLI_BINARY32 := kuskdcli-$(GOOS)_386
KUSKCLI_BINARY64 := kuskdcli-$(GOOS)_amd64

VERSION := $(shell awk -F= '/Version =/ {print $$2}' version/version.go | tr -d "\" ")

#MINER_RELEASE32 := miner-$(VERSION)-$(GOOS)_386
MINER_RELEASE64 := miner-$(VERSION)-$(GOOS)_amd64

#KUSKD_RELEASE32 := kuskd-$(VERSION)-$(GOOS)_386
KUSKD_RELEASE64 := kuskd-$(VERSION)-$(GOOS)_amd64

#KUSKCLI_RELEASE32 := kuskdcli-$(VERSION)-$(GOOS)_386
KUSKCLI_RELEASE64 := kuskdcli-$(VERSION)-$(GOOS)_amd64

#KUSK_RELEASE32 := kusk-$(VERSION)-$(GOOS)_386
KUSK_RELEASE64 := kusk-$(VERSION)-$(GOOS)_amd64

all: test target release-all install

kuskd:
	@echo "Building kuskd to cmd/kuskd/kuskd"
	@go build $(BUILD_FLAGS) -o cmd/kuskd/kuskd cmd/kuskd/main.go

kuskd-simd:
	@echo "Building SIMD version kuskd to cmd/kuskd/kuskd"
	@cd mining/tensority/cgo_algorithm/lib/ && make
	@go build -tags="simd" $(BUILD_FLAGS) -o cmd/kuskd/kuskd cmd/kuskd/main.go

kuskdcli:
	@echo "Building kuskdcli to cmd/kuskdcli/kuskdcli"
	@go build $(BUILD_FLAGS) -o cmd/kuskdcli/kuskdcli cmd/kuskdcli/main.go

install:
	@echo "Installing kuskd and kuskdcli to $(GOPATH)/bin"
	@go install ./cmd/kuskd
	@go install ./cmd/kuskdcli

target:
	mkdir -p $@

binary: target/$(KUSKD_BINARY64) target/$(KUSKCLI_BINARY64) target/$(MINER_BINARY64)

ifeq ($(GOOS),windows)
release: binary
	cd target && cp -f $(MINER_BINARY64) $(MINER_BINARY64).exe
	cd target && cp -f $(KUSKD_BINARY64) $(KUSKD_BINARY64).exe
	cd target && cp -f $(KUSKCLI_BINARY64) $(KUSKCLI_BINARY64).exe
	cd target && md5sum $(MINER_BINARY64).exe $(KUSKD_BINARY64).exe $(KUSKCLI_BINARY64).exe >$(KUSK_RELEASE64).md5
	cd target && zip $(KUSK_RELEASE64).zip $(MINER_BINARY64).exe $(KUSKD_BINARY64).exe $(KUSKCLI_BINARY64).exe $(KUSK_RELEASE64).md5
	cd target && rm -f $(MINER_BINARY64) $(KUSKD_BINARY64) $(KUSKCLI_BINARY64) $(MINER_BINARY64).exe $(KUSKD_BINARY64).exe $(KUSKCLI_BINARY64).exe $(KUSK_RELEASE64).md5
else
release: binary
	cd target && md5sum $(MINER_BINARY64) $(KUSKD_BINARY64) $(KUSKCLI_BINARY64) >$(KUSK_RELEASE64).md5
	cd target && tar -czf $(KUSK_RELEASE64).tgz $(MINER_BINARY64) $(KUSKD_BINARY64) $(KUSKCLI_BINARY64) $(KUSK_RELEASE64).md5
	cd target && rm -f $(MINER_BINARY64) $(KUSKD_BINARY64) $(KUSKCLI_BINARY64) $(KUSK_RELEASE64).md5
endif

release-all: clean
#	GOOS=darwin  make release
	GOOS=linux   make release
	GOOS=windows make release

clean:
	@echo "Cleaning binaries built..."
	@rm -rf cmd/kuskd/kuskd
	@rm -rf cmd/kuskdcli/kuskdcli
	@rm -rf cmd/miner/miner
	@rm -rf target
	@rm -rf $(GOPATH)/bin/kuskd
	@rm -rf $(GOPATH)/bin/kuskdcli
	@echo "Cleaning temp test data..."
	@rm -rf test/pseudo_hsm*
	@rm -rf blockchain/pseudohsm/testdata/pseudo/
	@echo "Cleaning sm2 pem files..."
	@rm -rf crypto/sm2/*.pem
	@echo "Done."

target/$(KUSKD_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/kuskd/main.go

target/$(KUSKCLI_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/kuskdcli/main.go

target/$(MINER_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/miner/main.go

test:
	@echo "====> Running go test"
	@go test -tags "network" $(PACKAGES)

benchmark:
	@go test -bench $(PACKAGES)

functional-tests:
	@go test -timeout=5m -tags="functional" ./test

ci: test functional-tests

.PHONY: all target release-all clean test benchmark
