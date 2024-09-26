# Image URL to use all building/pushing image targets
IMG ?= secret-reset:$(IMG_TAG)
IMG_REPO ?= ghcr.io/dana-team/$(NAME)
IMG_TAG ?= main
IMG_PULL_POLICY ?= Always
NAME ?= secret-reset
NAMESPACE ?= secret-reset-system

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test:  ## Run tests.
	go test $$(go list ./... | grep -v /e2e) -coverprofile cover.out

.PHONY: deploy
deploy: helm ## Deploy to the K8s cluster specified in ~/.kube/config.
	$(HELM) upgrade $(NAME) -n $(NAMESPACE) charts/$(NAME) --install --create-namespace \
		-f charts/$(NAME)/values.yaml \
		--set image.repository=$(IMG_REPO) \
		--set image.tag=$(IMG_TAG) \
		--set image.pullPolicy=$(IMG_PULL_POLICY) \
		--set config.env.AUTH_USERNAME=$(AUTH_USERNAME) \
		--set config.env.AUTH_CLIENT_SECRET=$(AUTH_CLIENT_SECRET) \
		--set config.env.AUTH_URL=$(AUTH_URL) \
		--set config.env.SECRET_NAME=$(SECRET_NAME) \
		--set config.env.SECRET_NAMESPACE=$(NAMESPACE)

.PHONY: undeploy
undeploy: helm ## Deploy to the K8s cluster specified in ~/.kube/config.
	$(HELM) uninstall $(NAME) -n $(NAMESPACE)

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

.PHONY: doc-chart
doc-chart: helm-docs helm
	$(HELM_DOCS) charts/

##@ Build

.PHONY: build
build: fmt vet ## Build manager binary.
	go build -o bin/manager cmd/main.go

.PHONY: run
run: fmt vet ## Run a controller from your host.
	go run ./cmd/main.go

# If you wish to build the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	$(CONTAINER_TOOL) build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${IMG}

##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
HELM ?= $(LOCALBIN)/helm
HELM_URL ?= https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
HELM_DOCS ?= $(LOCALBIN)/helm-docs-$(HELM_DOCS_VERSION)

## Tool Versions
GOLANGCI_LINT_VERSION ?= v1.60.3
HELM_DOCS_VERSION ?= v1.14.2

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,${GOLANGCI_LINT_VERSION})

.PHONY: helm
helm: $(HELM) ## Install helm on the local machine
$(HELM): $(LOCALBIN)
	wget -O $(LOCALBIN)/get-helm.sh $(HELM_URL)
	chmod 700 $(LOCALBIN)/get-helm.sh
	HELM_INSTALL_DIR=$(LOCALBIN) $(LOCALBIN)/get-helm.sh

.PHONY: helm-docs
helm-docs: $(HELM_DOCS)
$(HELM_DOCS): $(LOCALBIN)
	$(call go-install-tool,$(HELM_DOCS),github.com/norwoodj/helm-docs/cmd/helm-docs,$(HELM_DOCS_VERSION))

.PHONY: helm-plugins
helm-plugins: ## Install helm plugins on the local machine
	@if ! helm plugin list | grep -q 'diff'; then \
		helm plugin install https://github.com/databus23/helm-diff; \
	fi
	@if ! helm plugin list | grep -q 'git'; then \
		helm plugin install https://github.com/aslafy-z/helm-git; \
	fi
	@if ! helm plugin list | grep -q 's3'; then \
		helm plugin install https://github.com/hypnoglow/helm-s3; \
	fi
	@if ! helm plugin list | grep -q 'secrets'; then \
		helm plugin install https://github.com/jkroepke/helm-secrets; \
	fi

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary (ideally with version)
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f $(1) ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv "$$(echo "$(1)" | sed "s/-$(3)$$//")" $(1) ;\
}
endef