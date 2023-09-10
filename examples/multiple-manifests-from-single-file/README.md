<!-- omit in toc -->
# Example - Generate multiple manifests from a single file

- [Use case](#use-case)
- [Input](#input)
- [Manifest](#manifest)
- [Build](#build)
- [Output](#output)

## Use case

Currently, in Kustomize, you can use a single patch to override multiple resources. However, the vice versa
doesn't work in the same path (using a single resource as a template and generating multiple resources from it).

## Input

```yaml
# Base.
apiVersion: batch/v1
kind: CronJob
metadata:
  name: this-will-be-applied-to-multiple-cronjobs
spec:
  failedJobsHistoryLimit: 2
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: my-app
            image: my-app
            # The command will be updated via Merger.
            command: []
            imagePullPolicy: Always
          restartPolicy: Never
          imagePullSecrets:
          - name: docker-registry
  # The schedule will be updated via Merger.
  schedule: ""
  successfulJobsHistoryLimit: 0

```

```yaml
# Overlay 01.
apiVersion: batch/v1
kind: CronJob
metadata:
  name: cronjob01
spec:
  schedule: '0 1 * * *'
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: my-app
            command:
            - echo
            - cronjob01
```

```yaml
# Overlay 02.
apiVersion: batch/v1
kind: CronJob
metadata:
  name: cronjob02
spec:
  schedule: '0 2 * * *'
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: my-app
            command:
            - echo
            - cronjob02
```

```yaml
# Overlay 03.
apiVersion: batch/v1
kind: CronJob
metadata:
  name: cronjob03
spec:
  schedule: '0 3 * * *'
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: my-app
            command:
            - echo
            - cronjob03
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
    # exec:
    #   path: kustomize-plugin-merger
spec:
  resources:
  - name: my-cronjobs
    input:
      method: overlay
      files:
        # The same as in the KRM container above.
        root: /mnt
        sources:
        - input/cronjob-01.yaml
        - input/cronjob-02.yaml
        - input/cronjob-03.yaml
        destination: input/cronjob-common.yaml
    merge:
      strategy: combine
    output:
      format: raw
```

## Build

```shell
kustomize build  --enable-alpha-plugins --as-current-user .
```

## Output

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: cronjob01
spec:
  failedJobsHistoryLimit: 2
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - command:
            - echo
            - cronjob01
            image: my-app
            imagePullPolicy: Always
            name: my-app
          imagePullSecrets:
          - name: docker-registry
          restartPolicy: Never
  schedule: 0 1 * * *
  successfulJobsHistoryLimit: 0
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: cronjob02
spec:
  failedJobsHistoryLimit: 2
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - command:
            - echo
            - cronjob02
            image: my-app
            imagePullPolicy: Always
            name: my-app
          imagePullSecrets:
          - name: docker-registry
          restartPolicy: Never
  schedule: 0 2 * * *
  successfulJobsHistoryLimit: 0
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: cronjob03
spec:
  failedJobsHistoryLimit: 2
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - command:
            - echo
            - cronjob03
            image: my-app
            imagePullPolicy: Always
            name: my-app
          imagePullSecrets:
          - name: docker-registry
          restartPolicy: Never
  schedule: 0 3 * * *
  successfulJobsHistoryLimit: 0

```
