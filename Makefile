VERSION?=`git describe --tags`
DOCBIN?=mkdocs
BUILD_DATE := `date +%Y-%m-%d\ %H:%M`
VERSIONFILE := version.go

GOCMD=go
GOBUILD=$(GOCMD) build -mod=vendor -installsuffix cgo -a -x -ldflags "-w -v"
GOBUILD_DBG=$(GOCMD) build -x -mod=vendor
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=tts-cache
REPO=github.com/dodas-ts/docker-img_tts-go-cache

export GO111MODULE=on
# Force 64 bit architecture
export GOARCH=amd64

all: docker-img-build

build:
	env $(GOBUILD) -o $(BINARY_NAME)

build-debug:
	env CGO_ENABLED=0 $(GOBUILD_DBG) -o $(BINARY_NAME) -v

vendor:
	$(GOCMD) mod tidy
	$(GOCMD) vendor

docker-img-build: vendor build
	docker build . -t dodasts/tts-cache:$(VERSION)

windows-build:
	env GOOS=windows CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME).exe -v

macos-build:
	env GOOS=darwin CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME)_osx -v

gensrc:
	rm -f $(VERSIONFILE)
	@echo "package main" > $(VERSIONFILE)
	@echo "const (" >> $(VERSIONFILE)
	@echo "  VERSION = \"$(VERSION)\"" >> $(VERSIONFILE)
	@echo "  BUILD_DATE = \"$(BUILD_DATE)\"" >> $(VERSIONFILE)
	@echo ")" >> $(VERSIONFILE)

build-release: vendor gensrc docker-img-build

