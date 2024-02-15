COMMIT                := $(shell git rev-parse HEAD)
COMMIT_SHORT          := $(shell git rev-parse --short HEAD)
DATE                  := $(shell date +%Y-%m-%d)
BRANCH                := $(shell git rev-parse --abbrev-ref HEAD)
VERSION               ?= ${BRANCH}-${COMMIT_SHORT}
PKG_LDFLAGS           := github.com/prometheus/common/version
LDFLAGS               := -s -w -X ${PKG_LDFLAGS}.Version=${VERSION} -X ${PKG_LDFLAGS}.Revision=${COMMIT} -X ${PKG_LDFLAGS}.BuildDate=${DATE} -X ${PKG_LDFLAGS}.Branch=${BRANCH}

# Image URL to use all building/pushing image targets
DEFAULT_IMG_REPO := ghcr.io/spinkube/spin-operator
IMG ?= $(DEFAULT_IMG_REPO):$(shell git rev-parse --short HEAD)-dev

# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.28.3

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# CONTAINER_TOOL defines the container tool to be used for building images.
# we currently depend on Docker and `buildx` to ensure that we can build cross-arch
# images effectively. We may decide to change that in the future.
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

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./api/..." paths="./cmd/..." paths="./internal/..." paths="./pkg/..." output:crd:artifacts:config=config/crd/bases

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./api/..." paths="./cmd/..." paths="./internal/..." paths="./pkg/..."

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: manifests generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test ./... -coverprofile cover.out

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.54.2
golangci-lint:
	@[ -f $(GOLANGCI_LINT) ] || { \
	set -e ;\
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell dirname $(GOLANGCI_LINT)) $(GOLANGCI_LINT_VERSION) ;\
	}

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

.PHONY: helm-lint
helm-lint: helm-generate ## Lint the Helm chart
	$(HELM) lint ./charts/$(HELM_CHART)

.PHONY: lint-markdown
lint-markdown: ## Lint markdown files
	$(CONTAINER_TOOL) build --load -f format.Dockerfile -t markdown-formatter .
	$(CONTAINER_TOOL) run -e PRETTIER_MODE=check -v .:/usr/spin-operator markdown-formatter

.PHONY: lint-markdown-fix
lint-markdown-fix: ## Lint markdown files and perform fixes
	$(CONTAINER_TOOL) build --load -f format.Dockerfile -t markdown-formatter .
	$(CONTAINER_TOOL) run -e PRETTIER_MODE=write -v .:/usr/spin-operator markdown-formatter

##@ Build

.PHONY: golangci-build
golangci-build: ## Build manager binary.
	go build -ldflags "${LDFLAGS}" -a -o bin/manager cmd/main.go

.PHONY: build
build: manifests generate fmt vet golangci-build ## Build manager binary.

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run -ldflags "${LDFLAGS}" ./cmd/main.go

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	$(CONTAINER_TOOL) build --load -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${IMG}

PLATFORMS ?= linux/arm64,linux/amd64
.PHONY: docker-build-and-publish-all
docker-build-and-publish-all: ## Build the docker image for all supported platforms
	$(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag ${IMG} .

##@ Package

HELM_CHART := spin-operator

.PHONY: helm-generate
helm-generate: manifests kustomize helmify ## Create/update the Helm chart based on kustomize manifests. (Note: CRDs not included)
	$(KUSTOMIZE) build config/default | $(HELMIFY) -crd-dir -cert-manager-as-subchart -cert-manager-version v1.13.3 charts/$(HELM_CHART)
	rm -rf charts/$(HELM_CHART)/crds
	@# Copy the containerd-shim-spin SpinAppExecutor yaml from its canonical location into the chart
	cp config/samples/shim-executor.yaml charts/$(HELM_CHART)/templates/containerd-shim-spin-executor.yaml
	$(HELM) dep up charts/$(HELM_CHART)

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: manifests kustomize ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) apply -f -

.PHONY: uninstall
uninstall: manifests kustomize ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image $(DEFAULT_IMG_REPO)=${IMG}
	$(KUSTOMIZE) build config/default | $(KUBECTL) apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

HELM_RELEASE   ?= $(HELM_CHART)
HELM_NAMESPACE ?= $(HELM_CHART)
IMG_REPO := $(shell echo "${IMG}" | cut -d ':' -f 1)
IMG_TAG  := $(shell echo "${IMG}" | cut -d ':' -f 2)

.PHONY: helm-install
helm-install: helm-generate ## Install the Helm chart onto the K8s cluster specified in ~/.kube/config.
	$(HELM) upgrade --install \
		-n $(HELM_NAMESPACE) \
		--create-namespace \
		--set controllerManager.manager.image.repository=$(IMG_REPO) \
		--set controllerManager.manager.image.tag=$(IMG_TAG) \
		$(HELM_RELEASE) charts/$(HELM_CHART)

.PHONY: helm-upgrade
helm-upgrade: helm-install ## Upgrade the Helm release.

.PHONY: helm-uninstall
helm-uninstall: ## Delete the Helm release.
	$(HELM) delete \
		-n $(HELM_NAMESPACE) \
		$(HELM_RELEASE)

##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUBECTL ?= kubectl
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
HELM ?= helm
HELMIFY ?= $(LOCALBIN)/helmify

## Tool Versions
KUSTOMIZE_VERSION ?= v5.2.1
CONTROLLER_TOOLS_VERSION ?= v0.13.0
HELMIFY_VESRION ?= v0.4.10

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary. If wrong version is installed, it will be removed before downloading.
$(KUSTOMIZE): $(LOCALBIN)
	@if test -x $(LOCALBIN)/kustomize && ! $(LOCALBIN)/kustomize version | grep -q $(KUSTOMIZE_VERSION); then \
		echo "$(LOCALBIN)/kustomize version is not expected $(KUSTOMIZE_VERSION). Removing it before installing."; \
		rm -rf $(LOCALBIN)/kustomize; \
	fi
	test -s $(LOCALBIN)/kustomize || GOBIN=$(LOCALBIN) GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v5@$(KUSTOMIZE_VERSION)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary. If wrong version is installed, it will be overwritten.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen && $(LOCALBIN)/controller-gen --version | grep -q $(CONTROLLER_TOOLS_VERSION) || \
	GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
	$(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN)

.PHONY: helmify
helmify: $(HELMIFY) ## Download helmify locally if necessary.
$(HELMIFY): $(LOCALBIN)
	test -s $(LOCALBIN)/helmify || GOBIN=$(LOCALBIN) go install github.com/arttor/helmify/cmd/helmify@$(HELMIFY_VESRION)

.PHONY: container-registry-secret
container-registry-secret: ## Create a secret for the container registry
	@echo "Creating secret for container registry"
	@kubectl create ns $(HELM_NAMESPACE) || true
	@kubectl create secret docker-registry ghcr-credentials \
		--docker-server=https://ghcr.io \
		--docker-username=$(GH_USERNAME) \
		--docker-password=$(GH_PAT) \
		--namespace=$(HELM_NAMESPACE) || true

.PHONY: install-runtime-config
install-runtime-config: ## Install the runtime configuration for the operator
	@echo "Installing runtime configuration for the operator"
	kubectl apply -f spin-runtime-class.yaml

.PHONY: deploy-latest-cert-manager
deploy-latest-cert-manager: ## Deploy the latest version of cert-manager
	@echo "Deploying the latest version of cert-manager"
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.2/cert-manager.yaml

.PHONY: deploy-latest
deploy-latest: install deploy-latest-cert-manager container-registry-secret install-runtime-config 
	kubectl apply -n spin-operator -f config/samples/shim-executor.yaml
	$(MAKE) manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image $(DEFAULT_IMG_REPO)=ghcr.io/spinkube/spin-operator:20240215-045414-g0d84cbd
	$(KUSTOMIZE) build config/default | $(KUBECTL) apply -f -