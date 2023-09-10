[![CI](https://img.shields.io/github/actions/workflow/status/aabouzaid/kustomize-plugin-merger/.github%2Fworkflows%2Fgo-ci.yml?logo=github&label=CI&color=31c653)](https://github.com/aabouzaid/kustomize-plugin-merger/actions/workflows/go-ci.yml?query=branch%3Amain)
[![CodeQL](https://img.shields.io/github/actions/workflow/status/aabouzaid/kustomize-plugin-merger/.github%2Fworkflows%2Fgo-ci.yml?logo=github&label=CodeQL&color=31c653)](https://github.com/aabouzaid/kustomize-plugin-merger/actions/workflows/sec-codeql.yml?query=branch%3Amain)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/7815/badge)](https://www.bestpractices.dev/projects/7815)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/aabouzaid/kustomize-plugin-merger/badge)](https://securityscorecards.dev/viewer/?uri=github.com/aabouzaid/kustomize-plugin-merger)
[![Go Report Card](https://goreportcard.com/badge/github.com/aabouzaid/kustomize-plugin-merger)](https://goreportcard.com/report/github.com/aabouzaid/kustomize-plugin-merger)
[![GitHub Release](https://img.shields.io/github/v/release/aabouzaid/kustomize-plugin-merger?logo=github)](https://github.com/aabouzaid/kustomize-plugin-merger/releases)
[![Docker](https://img.shields.io/badge/Docker-available-blue?logo=docker&logoColor=white)](https://github.com/aabouzaid/kustomize-plugin-merger/pkgs/container/kustomize-generator-merger)
[![Go Reference](https://pkg.go.dev/badge/github.com/aabouzaid/kustomize-plugin-merger.svg)](https://pkg.go.dev/github.com/aabouzaid/kustomize-plugin-merger)
[![Renovatebot](https://img.shields.io/badge/Renovate-enabled-blue?logo=renovatebot)](https://github.com/aabouzaid/kustomize-plugin-merger/issues/7)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/aabouzaid/kustomize-plugin-merger/pulls)

<!-- omit in toc -->
# 🔀 Merger 🔀

A Kustomize generator plugin to merge YAML files seamlessly for real-world use cases.

**Merger** provides schemaless merge with different merge strategies (StrategicMerge).


- [Why](#why)
- [Features](#features)
- [Options](#options)
- [Common use cases](#common-use-cases)
  - [1. Generate multiple manifests from a single base](#1-generate-multiple-manifests-from-a-single-base)
  - [2. Merge lists in manifests without schema or a unique identifier](#2-merge-lists-in-manifests-without-schema-or-a-unique-identifier)
  - [3. Organize long manifests into smaller ones](#3-organize-long-manifests-into-smaller-ones)
- [TO-DO](#to-do)
- [Project status](#project-status)
- [Contributing](#contributing)
- [License](#license)


## Why

**The main goal of the `Merger` plugin is to provide a more flexible way to merge YAML files
and extend the strict merge capacities of Kustomize.**

Kustomize's default patch strategy for the lists (arrays) is `replace`, and to change that
for the Kubernetes `Custom Resource`, you must provide the OpenAPI schema of that custom resource,
which is only helpful if the lists of that CR have a unique id.

Many people are asking for such functionalities as [easy merging CRs](https://stackoverflow.com/q/73655002/4547221),
[creating a ConfigMap from arbitrary YAML files](https://stackoverflow.com/q/74547569/4547221),
and for more details on the challenge of providing OpenAPI schema to merge files, please read the following post:
[Set OpenAPI patch strategy for Kubernetes Custom Resources](https://tech.aabouzaid.com/2022/11/set-openapi-patch-strategy-for-kubernetes-custom-resources-kustomize.html).


## Features

- Generate multiple resources/manifests from a single base without copying the resources multiple times.
- Merge any manifests (even CustomResources) without needing their OpenAPI schema.
- Merge manifests with a list of maps without a unique identifier
  (when using `x-kubernetes-patch-merge-key` is not possible).
- Merge YAML files with different merge strategies (StrategicMerge).
- Merge applications configuration YAML files into a ConfigMap or Secret (WIP).


## Options

```yaml
---
apiVersion: generators.kustomize.aabouzaid.com/v1alpha1
kind: Merger
metadata:
  name: merge
    annotations:
      # Containerized KRM function.
      config.kubernetes.io/function: |
        container:
          image: ghcr.io/aabouzaid/kustomize-generator-merger
          mounts:
          - type: bind
            src: ./
            dst: /mnt
      # Exec KRM functions.
      # config.kubernetes.io/function: |
      # exec:
      #   path: kustomize-plugin-merger
spec:
  resources:
  - name: example
    input:
      # Available options: overlay,patch.
      # - Overlay: Produce multiple outputs by merging each source with the destination.
      # - Patch: Produce a single output by merging all sources together then with the destination.
      method: overlay
      files:
        # The same as in the KRM container above, omit it if Exec KRM is used.
        root: /mnt
        sources:
        - src01.yaml
        - src02.yaml
        destination: dst.yaml
    merge:
      # Available options: replace,append,combine.
      # - Replace: All keys in source will merge and replace what's in the destination.
      # - Append: Maps from source merged with destination, but the lists will be appended from source to destination.
      # - Combine: Maps from source merged with destination, but the lists will be combined together.
      strategy: combine
    output:
      # Available options: raw.
      format: raw
```


## Common use cases

This section shows a couple of use cases where Merger can help.

### 1. Generate multiple manifests from a single base

In this case, you have multiple `CronJobs`, all of which share the same body,
but each has a different command or other config.

[Use case full example](./examples/multiple-manifests-from-single-file/README.md).

### 2. Merge lists in manifests without schema or a unique identifier

Currently, in Kustomize, it's not possible to merge resources without a unique identifier, even with Open API schema.

It's possible to do that using the merge strategy `append` in Merger (later on, `combineWithKey` will also be supported).

[Use case full example](./examples/manifest-lists-without-schema/README.md).

### 3. Organize long manifests into smaller ones

In some use cases (e.g., [Crossplane Compositions](https://docs.crossplane.io/latest/concepts/compositions/)),
you could have a really long YAML manifest, and it's hard to read. You can split that file
and use the Merger `patch` input method to make it a single manifest again.

[Use case full example](./examples/long-omni-manifest/README.md).


## TO-DO

- Support `ConfigMap` or `Secret` as an output.
- Support `combine` merge strategy with an identifier key (similar to `x-kubernetes-patch-merge-key`).
- Provide better docs for Merger options.


## Project status

Please note that this project is still under development and could have breaking changes,
but it will follow the SemVer convention.


## Contributing

All feedback and contributions are welcome. If you find any issues or want to contribute,
please feel free to [fill an issue](https://github.com/aabouzaid/kustomize-plugin-merger/issues)
or [create a PR](https://github.com/aabouzaid/kustomize-plugin-merger/pulls).


## License

Merger is an open-source software licensed using the Apache License 2.0. Please see [LICENSE](LICENSE) for details.
