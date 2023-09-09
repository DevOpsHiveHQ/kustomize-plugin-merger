[![CI](https://img.shields.io/github/actions/workflow/status/aabouzaid/kustomize-plugin-merger/.github%2Fworkflows%2Fgo-ci.yml?logo=github&label=CI&color=31c653)](https://github.com/aabouzaid/kustomize-plugin-merger/actions/workflows/go-ci.yml?query=branch%3Amain)
[![CodeQL](https://img.shields.io/github/actions/workflow/status/aabouzaid/kustomize-plugin-merger/.github%2Fworkflows%2Fgo-ci.yml?logo=github&label=CodeQL&color=31c653)](https://github.com/aabouzaid/kustomize-plugin-merger/actions/workflows/sec-codeql.yml?query=branch%3Amain)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/aabouzaid/kustomize-plugin-merger/badge)](https://securityscorecards.dev/viewer/?uri=github.com/aabouzaid/kustomize-plugin-merger)
[![go.mod](https://img.shields.io/github/go-mod/go-version/aabouzaid/kustomize-plugin-merger?logo=go&logoColor=white)](go.mod)
[![Renovatebot](https://img.shields.io/badge/Renovate-enabled-blue?logo=renovatebot)](https://github.com/aabouzaid/kustomize-plugin-merger/issues/7)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/aabouzaid/kustomize-plugin-merger/pulls)

# Merger (WIP)

A Kustomize generator plugin to merge YAML files seamlessly (schemaless merge).


## Why

TBA


## Use Cases

TBA


## Build

```sh
go build -o 'dist/merger' .
```


## Example

```yaml
---
apiVersion: generators.kustomize.aabouzaid.com/v1alpha1
kind: Merger
metadata:
  name: merge
spec:
  resources:
  - name: my-envs
    input:
      # Available options: overlay,patch.
      method: overlay
      files:
        sources:
        - example/input/dev.yaml
        - example/input/stage.yaml
        destination: example/input/common.yaml
    merge:
      # Available options: replace,append,combine.
      strategy: combine
    output:
      # Available options: raw.
      format: raw
```
