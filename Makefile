FORCE: ;

SELF := $(MAKE)
OS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH ?= $(shell uname -m | sed 's/x86_64/amd64/g')
GO := $(shell which go 2>/dev/null)
MOCKGEN := $(shell which mockgen)

INSTALL_DIR ?= /usr/local/sbin
RELEASE_DIR ?= release
APP ?= $(shell basename `pwd`)

null  :=
space := $(null) #
comma := ,
PACKAGES := $(shell pwd)/cmd/api,$(shell pwd)/cmd/server
PKGS ?= $(subst $(comma),$(space),$(PACKAGES))

export GCARCH_POSTGRES_HOST=localhost
export GCARCH_POSTGRES_USER=postgres
export GCARCH_POSTGRES_PASSWORD=postgres
export GCARCH_POSTGRES_DB=gcarch_example
export GCARCH_PROMETHEUS_PUSHGATEWAY=http://localhost:9091
export GCARCH_API_PORT=9000

## Run the api
go/run/api:
	$(GO) run cmd/api/main.go

## Build binary for all platforms
go/build: $(GO)
	$(call assert-set,GO)
ifeq ($(RELEASE_ARCH),)
	gox -output "${RELEASE_DIR}/${APP}_{{.Dir}}_{{.OS}}_{{.Arch}}" $(PKGS)
else
	gox -osarch="$(RELEASE_ARCH)" -output "${RELEASE_DIR}/${APP}_{{.Dir}}_{{.OS}}_{{.Arch}}" $(PKGS)
endif

go/watch:
	air -c .air.toml

go/exec/api:
	${RELEASE_DIR}/${APP}_api_${OS}_${ARCH}

## Install dependencies
go/dep:
	$(call assert-set,GO)
	$(GO) mod vendor

## Install development dependencies
go/dep/dev: $(GO)
	$(call assert-set,GO)
	$(GO) get -u -v github.com/golang/mock/gomock
	$(GO) get -u -v github.com/golang/mock/mockgen
	$(GO) get -u -v github.com/golang/mock/mockgen/model
	$(GO) get -u -v github.com/mitchellh/gox
	$(GO) get -u -v github.com/cosmtrek/air
	@$(SELF) go/dep

## Clean compiled binary
go/clean:
	rm -rf $(RELEASE_DIR)

## Clean compiled binary and dependency
go/clean/all: go/clean
	rm -rf vendor

## Install cli
go/install: $(APP) go/build
	cp $(RELEASE_DIR)/$(APP) $(INSTALL_DIR)
	chmod 555 $(INSTALL_DIR)/$(APP)

## Lint code
go/lint: $(GO) go/vet
	$(call assert-set,GO)
	find . ! -path "*/vendor/*" -type f -name '*.go' | xargs -n 1 golint

## Vet code
go/vet: $(GO)
	$(call assert-set,GO)
	find . ! -path "*/vendor/*" ! -path "*/.glide/*" -type f -name '*.go' | xargs $(GO) tool vet -v

## Format code according to Golang convention
go/fmt: $(GO)
	$(call assert-set,GO)
	find . ! -path "*/vendor/*" ! -path "*/.glide/*" -type f -name '*.go' | xargs -n 1 gofmt -w -l -s

## Run tests
go/test: $(GO)
	$(call assert-set,GO)
ifneq ($(LINT),true)
	$(GO) test $(shell $(GO) list ./... | grep -v /vendor/) -coverprofile=./ops/coverage.out
endif

## Open coverage
go/coverage/html: $(GO)
	$(call assert-set,GO)
	$(GO) tool cover -html=./ops/coverage.out

## Show coverage %
go/coverage/func: $(GO)
	$(call assert-set,GO)
	$(GO) tool cover -func=./ops/coverage.out

## Run docker compose
compose/up:
	docker-compose -f ./ops/docker-compose.yml up -d