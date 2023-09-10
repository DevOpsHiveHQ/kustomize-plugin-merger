CONTROLLER_GEN_HOME := /tmp/controller-gen
CONTROLLER_GEN_VERSION := v0.13.0

#
# Schema targets.
#

.PHONY: schema.install-tools 
schema.install-tools:
		mkdir -p $(CONTROLLER_GEN_HOME);\
		export GOBIN=$(CONTROLLER_GEN_HOME);\
		go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION);\

.PHONY: schema.generate
schema.generate:
	$(CONTROLLER_GEN_HOME)/controller-gen \
		crd:ignoreUnexportedFields=true,headerFile=pkg/merger/schema/header.yaml \
		paths=./... output:crd:dir=pkg/merger/schema; \
	sed -e 's/CustomResourceDefinition/KRMFunctionDefinition/g' \
    -e 's|apiextensions.k8s.io/v1|config.kubernetes.io/v1alpha1|g' \
    -i  pkg/merger/schema/generators.kustomize.aabouzaid.com_mergers.yaml

#
# Golang targets.
#

.PHONY: go.build
go.build:
	go build -o 'dist' .

.PHONY: go.format
go.format:
	gofumpt -l -w .
