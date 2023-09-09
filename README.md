[![Go CI](https://github.com/aabouzaid/kustomize-plugin-merger/actions/workflows/go-ci.yml/badge.svg)](https://github.com/aabouzaid/kustomize-plugin-merger/actions/workflows/go-ci.yml)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/aabouzaid/kustomize-plugin-merger/badge)](https://securityscorecards.dev/viewer/?uri=github.com/aabouzaid/kustomize-plugin-merger)
[![go.mod](https://img.shields.io/github/go-mod/go-version/aabouzaid/kustomize-plugin-merger)](go.mod)
[![License](https://img.shields.io/github/license/aabouzaid/kustomize-plugin-merger)](https://github.com/aabouzaid/kustomize-plugin-merger/blob/main/LICENSE)

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
