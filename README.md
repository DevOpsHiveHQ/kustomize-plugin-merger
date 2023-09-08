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
