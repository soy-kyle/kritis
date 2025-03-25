# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PWD = $(shell pwd)
GOOS ?= $(shell go env GOOS)
GOARCH = amd64
BUILD_DIR ?= ./out
COMMIT ?= $(shell git rev-parse HEAD)
VERSION ?= v0.2.2
IMAGE_TAG ?= $(COMMIT)

# Used for integration testing. example:
# "make -e GCP_PROJECT=kritis-int integration-local"
GCP_PROJECT ?= PLEASE_SET_GCP_PROJECT
GCP_ZONE ?= us-central1-a
GCP_CLUSTER ?= kritis-integration-test

%.exe: %
	mv $< $@

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

GITHUB_ORG := github.com/soy-kyle
GITHUB_PROJECT := kritis
REPOPATH ?= $(GITHUB_ORG)/$(GITHUB_PROJECT)
RESOLVE_TAGS_PROJECT := resolve-tags

SUPPORTED_PLATFORMS := linux-$(GOARCH) darwin-$(GOARCH) windows-$(GOARCH).exe
RESOLVE_TAGS_PATH = cmd/kritis/kubectl/plugins/resolve
RESOLVE_TAGS_PACKAGE = $(REPOPATH)/$(RESOLVE_TAGS_PATH)
RESOLVE_TAGS_KUBECTL_DIR = ~/.kube/plugins/resolve-tags

# Path to credentials. Must have a basename of gac.json.
GAC_CREDENTIALS_PATH ?= $(CURDIR)/.secrets/$(GCP_PROJECT)/gac.json

.PHONY: test
test: cross
	./hack/check-fmt.sh
	./hack/boilerplate.sh
	./hack/verify-codegen.sh
	./hack/dependencies.sh
	./hack/test.sh

GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GO_LD_RESOLVE_FLAGS :=""
GO_BUILD_TAGS := ""

.PRECIOUS: $(foreach platform, $(SUPPORTED_PLATFORMS), $(BUILD_DIR)/$(GITHUB_PROJECT)-$(platform))

$(BUILD_DIR)/$(RESOLVE_TAGS_PROJECT): $(BUILD_DIR)/$(RESOLVE_TAGS_PROJECT)-$(GOOS)-$(GOARCH)
	cp $(BUILD_DIR)/$(RESOLVE_TAGS_PROJECT)-$(GOOS)-$(GOARCH) $@

.PHONY: cross
cross: $(foreach platform, $(SUPPORTED_PLATFORMS), $(BUILD_DIR)/$(RESOLVE_TAGS_PROJECT)-$(platform))

$(BUILD_DIR)/$(RESOLVE_TAGS_PROJECT)-%-$(GOARCH): $(GO_FILES) $(BUILD_DIR)
	GOOS=$* GOARCH=$(GOARCH) CGO_ENABLED=0 go build -ldflags $(GO_LD_RESOLVE_FLAGS) -tags $(GO_BUILD_TAGS) -o $@ $(RESOLVE_TAGS_PACKAGE)

.PHONY: cross-tar
cross-tar: $(foreach platform, $(SUPPORTED_PLATFORMS), $(BUILD_DIR)/$(RESOLVE_TAGS_PROJECT)-$(platform).tar.gz)

$(BUILD_DIR)/$(RESOLVE_TAGS_PROJECT)-%.tar.gz: cross
	tar -czf $@ -C $(RESOLVE_TAGS_PATH) plugin.yaml -C $(PWD)/out/ resolve-tags-$*

.PHONY: install-plugin
install-plugin: $BUILD_DIR)/$(RESOLVE_TAGS_PROJECT)
	mkdir -p $(RESOLVE_TAGS_KUBECTL_DIR)
	cp $(BUILD_DIR)/$(RESOLVE_TAGS_PROJECT) $(RESOLVE_TAGS_KUBECTL_DIR)
	cp cmd/kritis/kubectl/plugins/resolve/plugin.yaml $(RESOLVE_TAGS_KUBECTL_DIR)

GO_LDFLAGS := -extldflags "-static"
GO_LDFLAGS += -X github.com/soy-kyle/kritis/cmd/kritis/version.Commit=$(COMMIT)
GO_LDFLAGS += -X github.com/soy-kyle/kritis/cmd/kritis/version.Version=$(VERSION)
GO_LDFLAGS += -w -s # Drop debugging symbols.

REGISTRY?=gcr.io/kritis-project
# TODO(tstromberg): Determine if it is possible to combine these two variables.
TEST_REGISTRY?=gcr.io/$(GCP_PROJECT)
SERVICE_PACKAGE = $(REPOPATH)/cmd/kritis/admission
SIGNER_PACKAGE = $(REPOPATH)/cmd/kritis/signer
GCR_SIGNER_PACKAGE = $(REPOPATH)/cmd/kritis/gcr-signer


out/kritis-server: $(GO_FILES)
	GOARCH=$(GOARCH) GOOS=linux CGO_ENABLED=0 go build -ldflags "$(GO_LDFLAGS)" -o $@ $(SERVICE_PACKAGE)

out/signer: $(GO_FILES)
	GOARCH=$(GOARCH) GOOS=linux CGO_ENABLED=0 go build -ldflags "$(GO_LDFLAGS)" -o $@ $(SIGNER_PACKAGE)

out/gcr-signer: $(GO_FILES)
	GOARCH=$(GOARCH) GOOS=linux CGO_ENABLED=0 go build -ldflags "$(GO_LDFLAGS)" -o $@ $(GCR_SIGNER_PACKAGE)

.PHONY: build-image
build-image: out/kritis-server
	docker build -t $(REGISTRY)/kritis-server:$(IMAGE_TAG) -f deploy/Dockerfile .

# build-test-image locally builds images for use in integration testing.
.PHONY: build-test-image
build-test-image: out/kritis-server
	docker build -t $(TEST_REGISTRY)/kritis-server:$(IMAGE_TAG) -f deploy/Dockerfile .

.PHONY: signer-image
signer-image: out/signer
	docker build -t $(REGISTRY)/kritis-signer:$(IMAGE_TAG) -f deploy/kritis-signer/Dockerfile .

.PHONY: gcr-signer-image
gcr-signer-image: out/gcr-signer
	docker build -t $(REGISTRY)/gcr-kritis-signer:$(IMAGE_TAG) -f deploy/gcr-kritis-signer/Dockerfile .

.PHONY: signer-push-image
signer-push-image: signer-image
	docker push $(REGISTRY)/kritis-signer:$(IMAGE_TAG)

.PHONY: gcr-signer-push-image
gcr-signer-push-image: gcr-signer-image
	docker push $(REGISTRY)/gcr-kritis-signer:$(IMAGE_TAG)

HELM_HOOKS = preinstall postinstall predelete

$(HELM_HOOKS): $(GO_FILES)
	GOARCH=$(GOARCH) GOOS=linux CGO_ENABLED=0 go build -ldflags "$(GO_LDFLAGS)" -o out/$@ $(REPOPATH)/helm-hooks/$@

.PHONY: %-test-image
%-test-image: $(HELM_HOOKS)
	docker build -t $(TEST_REGISTRY)/$*:$(IMAGE_TAG) -f helm-hooks/Dockerfile . --build-arg stage=$*

.PHONY: %-image
%-image: $(HELM_HOOKS)
	docker build -t $(REGISTRY)/$*:$(IMAGE_TAG) -f helm-hooks/Dockerfile . --build-arg stage=$*

.PHONY: helm-release-image
helm-release-image:
	docker build -t $(REGISTRY)/helm-release:$(IMAGE_TAG) -f helm-release/Dockerfile .

.PHONY: helm-install-from-head
helm-install-from-head:
	helm install --set repository=$(TEST_REGISTRY)/ --set image.tag=$(COMMIT) ./kritis-charts

clean:
	rm -rf $(BUILD_DIR)

.PHONY: integration-prod
integration-prod: cross
	go test -ldflags "$(GO_LDFLAGS)" -v -tags integration \
		$(REPOPATH)/integration \
		-timeout 15m \
		-gac-credentials=/tmp/gac.json \
		-gcp-project=kritis-int-test \
		-gke-cluster-name=test-cluster-2 $(EXTRA_TEST_FLAGS)

.PHONY: build-push-image
build-push-image: build-image preinstall-image postinstall-image predelete-image
	docker push $(REGISTRY)/kritis-server:$(IMAGE_TAG)
	docker push $(REGISTRY)/preinstall:$(IMAGE_TAG)
	docker push $(REGISTRY)/postinstall:$(IMAGE_TAG)
	docker push $(REGISTRY)/predelete:$(IMAGE_TAG)

.PHONY: build-push-test-image
build-push-test-image: build-test-image preinstall-test-image postinstall-test-image predelete-test-image
	docker push $(TEST_REGISTRY)/kritis-server:$(IMAGE_TAG)
	docker push $(TEST_REGISTRY)/preinstall:$(IMAGE_TAG)
	docker push $(TEST_REGISTRY)/postinstall:$(IMAGE_TAG)
	docker push $(TEST_REGISTRY)/predelete:$(IMAGE_TAG)

.PHONY: integration-in-docker
integration-in-docker: build-push-image
	docker build \
		-f deploy/$(GCP_PROJECT)/Dockerfile \
		--target integration \
		-t $(REGISTRY)/kritis-integration:$(IMAGE_TAG) .
	docker run \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(HOME)/tmp:/tmp \
		-v $(HOME)/.config/gcloud:/root/.config/gcloud \
		-v $(GOOGLE_APPLICATION_CREDENTIALS):$(GOOGLE_APPLICATION_CREDENTIALS) \
		-e REMOTE_INTEGRATION=true \
		-e DOCKER_CONFIG=/root/.docker \
		-e GOOGLE_APPLICATION_CREDENTIALS=$(GOOGLE_APPLICATION_CREDENTIALS) \
		$(REGISTRY)/kritis-integration:$(IMAGE_TAG)


# Fully setup local integration testing: only needs to run just once
# TODO: move entire setup into bash script
.PHONY: setup-integration-local
setup-integration-local:
	gcloud --project=$(GCP_PROJECT) services enable container.googleapis.com
	gcloud --project=$(GCP_PROJECT) container clusters describe $(GCP_CLUSTER) >/dev/null \
		|| gcloud --project=$(GCP_PROJECT) container clusters create $(GCP_CLUSTER) \
		--num-nodes=2 --zone=$(GCP_ZONE)
	gcloud --project=$(GCP_PROJECT) container clusters get-credentials $(GCP_CLUSTER)
	mkdir -p $(dir $(GAC_CREDENTIALS_PATH))
	test -n "$$(gcloud --project=$(GCP_PROJECT)  iam service-accounts list --filter "email: kritis-ca-admin@" --format "value(email)")" || \
	gcloud --project=$(GCP_PROJECT) iam service-accounts create kritis-ca-admin
	test -s $(GAC_CREDENTIALS_PATH) \
		|| gcloud --project=$(GCP_PROJECT) iam service-accounts keys \
		create $(GAC_CREDENTIALS_PATH) --iam-account kritis-ca-admin@${GCP_PROJECT}.iam.gserviceaccount.com
	kubectl get serviceaccount --namespace kube-system tiller >/dev/null 2>&1 || \
	    kubectl create serviceaccount --namespace kube-system tiller
	kubectl get clusterrolebinding tiller-cluster-rule >/dev/null 2>&1 || \
	kubectl create clusterrolebinding tiller-cluster-rule \
		  --clusterrole=cluster-admin \
		    --serviceaccount=kube-system:tiller
	helm init --wait --service-account tiller
	gcloud --project=$(GCP_PROJECT) services enable containerregistry.googleapis.com
	gcloud --project=$(GCP_PROJECT) services enable containeranalysis.googleapis.com
	gcloud --project=$(GCP_PROJECT) services enable containerscanning.googleapis.com
	gcloud -q container images add-tag \
		gcr.io/kritis-tutorial/acceptable-vulnz@sha256:2a81797428f5cab4592ac423dc3049050b28ffbaa3dd11000da942320f9979b6 \
		gcr.io/$(GCP_PROJECT)/acceptable-vulnz:latest
	gcloud -q container images add-tag \
		gcr.io/kritis-tutorial/java-with-vulnz@sha256:358687cfd3ec8e1dfeb2bf51b5110e4e16f6df71f64fba01986f720b2fcba68a \
		gcr.io/$(GCP_PROJECT)/java-with-vulnz:latest
	gcloud -q container images add-tag \
		gcr.io/kritis-tutorial/nginx-digest-whitelist:latest \
		gcr.io/$(GCP_PROJECT)/nginx-digest-whitelist:latest
	gcloud -q container images add-tag \
		gcr.io/kritis-tutorial/nginx-no-digest-breakglass:latest \
		gcr.io/$(GCP_PROJECT)/nginx-no-digest-breakglass:latest
	gcloud -q container images add-tag \
		gcr.io/kritis-tutorial/nginx-no-digest:latest \
		gcr.io/$(GCP_PROJECT)/nginx-no-digest:latest
	gcloud projects add-iam-policy-binding ${GCP_PROJECT} \
		--member=serviceAccount:kritis-ca-admin@${GCP_PROJECT}.iam.gserviceaccount.com \
		--role=roles/containeranalysis.notes.occurrences.viewer
	gcloud projects add-iam-policy-binding ${GCP_PROJECT} \
		--member=serviceAccount:kritis-ca-admin@${GCP_PROJECT}.iam.gserviceaccount.com \
		--role=roles/containeranalysis.occurrences.viewer
	gcloud projects add-iam-policy-binding ${GCP_PROJECT} \
		--member=serviceAccount:kritis-ca-admin@${GCP_PROJECT}.iam.gserviceaccount.com \
		--role=roles/containeranalysis.occurrences.editor
	./hack/setup-containeranalysis-resources.sh --project $(GCP_PROJECT)

# Fully clean-up local integration testing resources
.PHONY: clean-integration-local
clean-integration-local:
	gcloud --project=$(GCP_PROJECT) container clusters describe $(GCP_CLUSTER) >/dev/null \
		&& gcloud --project=$(GCP_PROJECT) container clusters delete $(GCP_CLUSTER)

# Just run the integration tests, assuming setup is done and test image is updated.
# make -e GCP_PROJECT=${PROJECT} just-the-integration-test
.PHONY: just-the-integration-test
just-the-integration-test:
	echo "Test cluster: $(GCP_CLUSTER) Test project: $(GCP_PROJECT)"
	go test -ldflags "$(GO_LDFLAGS)" -v -tags integration \
		$(REPOPATH)/integration \
		-timeout 30m \
		-gac-credentials=$(GAC_CREDENTIALS_PATH) \
		-gcp-project=$(GCP_PROJECT) \
		-gke-zone=$(GCP_ZONE) \
		-gke-cluster-name=$(GCP_CLUSTER) $(EXTRA_TEST_FLAGS)

# integration-local requires that "setup-integration-local" has been run at least once.
#
# Example usage, to run a single test without cleaning up:
#
#  make -e GCP_PROJECT=my-project \
#    EXTRA_TEST_FLAGS="-run TestKritisISPLogic/vulnz/acceptable-vulnz-replicaset.yaml --cleanup=false" \
#    integration-local
.PHONY: integration-local
integration-local: build-push-test-image just-the-integration-test
