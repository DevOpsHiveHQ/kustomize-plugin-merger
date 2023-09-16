<!-- omit in toc -->
# Example - Merge non-manifest files and store them into ConfigMap or Secret

- [Use case](#use-case)
- [Input](#input)
- [Manifest](#manifest)
- [Build](#build)
- [Output](#output)

## Use case

No plans from Kustomize to support non-manifest files merge and storing them into ConfigMap
or Secret. Using Merger you can merge any YAML files like application configuration.

In the example, two Prometheus config are merged together and stored as ConfigMap.

## Input

```yaml
# prometheus-base.yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets:
            - alertmanager:9093

# This will be added via Merger.
rule_files:

scrape_configs:
  - job_name: "prometheus"
    static_configs:
      - targets: ["localhost:9090"]
```

```yaml
# prometheus.yaml
rule_files:
  - "first_rules.yml"
  - "second_rules.yml"
```

## Manifest

```yaml
# kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

generators:
- merger.yaml
```

```yaml
# merger.yaml
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
    #   exec:
    #     path: kustomize-plugin-merger
spec:
  resources:
  - name: my-prometheus
    input:
      method: overlay
      files:
        # The same as in the KRM container above, omit it if Exec KRM is used.
        root: /mnt
        sources:
        - input/prometheus.yaml
        destination: input/prometheus-base.yaml
    merge:
      strategy: combine
    output:
      format: configmap
```

## Build

```shell
kustomize build  --enable-alpha-plugins --as-current-user .
```

## Output

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-prometheus
data:
  prometheus.yaml: |
    alerting:
        alertmanagers:
            - static_configs:
                - targets:
                    - alertmanager:9093
    global:
        evaluation_interval: 15s
        scrape_interval: 15s
    rule_files:
        - first_rules.yml
        - second_rules.yml
    scrape_configs:
        - job_name: prometheus
          static_configs:
            - targets:
                - localhost:9090
```