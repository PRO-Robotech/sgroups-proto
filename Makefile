SHELL:=/bin/sh
export GOSUMDB=off
export GO111MODULE=on

$(value $(shell [ ! -d "$(CURDIR)/bin" ] && mkdir -p "$(CURDIR)/bin"))
export GOBIN=$(CURDIR)/bin
GO?=$(shell which go)

platform?=$(shell uname -s)/$(shell uname -m)
os?=$(strip $(filter Linux Darwin,$(word 1,$(subst /, ,$(platform)))))
arch?=$(strip $(filter x86_64 arm64,$(word 2,$(subst /, ,$(platform)))))

.PHONY: help
help: ##display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


# install project dependencies
.PHONY: go-deps
go-deps: ##install golang dependencies
	@echo Check go modules dependencies... && \
	$(GO) mod tidy && go mod vendor && go mod verify && \
	echo -=OK=-


BUF_VERSION:=v1.64.0
BUF_REPO:=https://github.com/bufbuild/buf/releases/download/${BUF_VERSION}/buf-$(os)-$(arch)

BUF:=$(GOBIN)/buf

.PHONY: .install-buf
.install-buf:
ifneq ($(wildcard $(BUF)),)
	@echo >/dev/null
else
ifeq ($(filter Linux Darwin,$(os)),)
	$(error os=$(os) but must be in [Linux|Darwin])
endif
ifeq ($(filter x86_64 arm64,$(arch)),)
	$(error arch=$(arch) but must be in [x86_64|arm64])
endif
	@echo "Downloading $(BUF_REPO)..."
	@mkdir -p $(GOBIN)
	@wget -q -O $(GOBIN)/buf $(BUF_REPO)
	@chmod +x $(GOBIN)/buf
	@echo "buf installed to $(GOBIN)/buf"
endif

.PHONY: generate-api
generate-api: | proto-deps ##generate API code
	@(\
	dest=$(CURDIR)/pkg/api && \
	rm -rf $$dest 2>/dev/null && \
	mkdir -p $$dest && \
	echo generating API in \"$$dest\" ... && \
	$(BUF) generate api --template api/buf.gen.yaml \
		--path api/common \
		--path api/sgroups \
		--exclude-path api/sgroups/v1/queries.proto &&\
	$(MAKE) go-deps && \
	echo -=OK=- ;\
	)

.PHONY: proto-lint
proto-lint: | .install-buf ##lint protos
	@echo Running lint proto... && \
	$(BUF) lint api && \
	echo -=OK=-

.PHONY: proto-validate
proto-validate: | .install-buf ##validate protos
	@echo "Running validation..." && \
	$(BUF) build api && \
	echo -=OK=-

.PHONY: proto-format
proto-format: | .install-buf ##format protos
	@echo Formatting protos with buf... && \
	$(BUF) format -w api --exclude-path api/swagger-title.proto && \
	echo -=OK=-

.PHONY: proto-format-check
proto-format-check: | .install-buf ##check proto formatting
	@echo Checking proto formatting... && \
	$(BUF) format --diff --exit-code api --exclude-path api/swagger-title.proto && \
	echo -=OK=-

.PHONY: proto-deps
proto-deps: | .install-buf ##update proto dependencies
	@echo "Updating proto dependencies..." && \
	$(BUF) dep update api && \
	echo -=OK=-


.PHONY: test
test: ##run tests
	@echo Running tests... && \
	$(GO) clean -testcache && go test -v ./... && \
	echo -=OK=-


