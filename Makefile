# env defines
GOOS=$(shell go env GOOS)
ARCH=$(shell arch)
VERSION=$(shell cat ./VERSION)
GO_VERSION=$(shell go env GOVERSION)
GIT_COMMIT_ID=$(shell git rev-parse HEAD)
GIT_DESCRIBE=$(shell git describe --always)
OS=$(if $(GOOS),$(GOOS),linux)

# go command defines
GO_BUILD=go build
GO_MOD_TIDY=$(go mod tidy -compat 1.19)
GO_BUILD_WITH_INFO=$(GO_BUILD) -ldflags "\
	-X 'preinstall/defs/compiledef._appVersion=$(VERSION)' \
	-X 'preinstall/defs/compiledef._goVersion=$(GO_VERSION)'\
	-X 'preinstall/defs/compiledef._gitCommitID=$(GIT_COMMIT_ID)'\
	-X 'preinstall/defs/compiledef._gitDescribe=$(GIT_DESCRIBE)'"

# package defines
PKG_PERFIX=yashandb-preinstall-$(VERSION)
PKG=$(PKG_PERFIX)-$(OS)-$(ARCH).tar.gz

# build defines
BUILD_PATH=./build
PKG_PATH=$(BUILD_PATH)/$(PKG_PERFIX)
BIN_PATH=$(PKG_PATH)/bin
LOG_PATH=$(PKG_PATH)/log
RESULT_PATH=$(PKG_PATH)/result
DOCS_PATH=$(PKG_PATH)/docs
PLUGINS_PATH=$(PKG_PATH)/plugins
BIN_PREINSTALL=$(BUILD_PATH)/preinstall
BIN_FILES=$(BIN_PREINSTALL)
DIR_TO_MAKE=$(BIN_PATH) $(LOG_PATH) $(RESULT_PATH) $(DOCS_PATH) $(PLUGINS_PATH)
FILE_TO_COPY=./config


.PHONY: clean force go_build

build: pre_build go_build
	@./.resolve-goimports.sh -q
	@mv $(BIN_FILES) $(BIN_PATH)
	@> $(LOG_PATH)/preinstall.log
	@cd $(PKG_PATH);ln -s ./bin/preinstall ./preinstall
	@cd $(BUILD_PATH);tar -cvzf $(PKG) $(PKG_PERFIX)/

clean:
	rm -rf $(BUILD_PATH)

go_build: 
	$(GO_MOD_TIDY)
	$(GO_BUILD_WITH_INFO) -o $(BIN_PREINSTALL) ./cmd/*.go

pre_build:
	@mkdir -p $(DIR_TO_MAKE) 
	@cp -r $(FILE_TO_COPY) $(PKG_PATH)
	@cp ./plugins/$(ARCH)/* $(PLUGINS_PATH)
	@cp ./README.md $(DOCS_PATH)/preinstall.md

force: clean build