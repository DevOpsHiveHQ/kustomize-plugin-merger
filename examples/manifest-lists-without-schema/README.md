<!-- omit in toc -->
# Example - Append CustomResources lists without schema

- [Use case](#use-case)
- [Input](#input)
- [Manifest](#manifest)
- [Build](#build)
- [Output](#output)

## Use case

Currently, in Kustomize, it's not possible to append CustomResources without providing an Open API schema
(which is a lot of hassles for many use cases, to be honest).

## Input

```yaml
# Base.
apiVersion: iam.cnrm.cloud.google.com/v1beta1
kind: IAMPolicy
metadata:
  name: storage-admin-policy
spec:
  resourceRef:
    apiVersion: iam.cnrm.cloud.google.com/v1beta1
    kind: IAMServiceAccount
    name: storage-admin
    namespace: cnrm-system
  bindings:
    - role: roles/iam.workloadIdentityUser
      members:
        - serviceAccount:mysten-sui.svc.id.goog[monitoring/thanos-compactor]
```

```yaml
# Overlay.
apiVersion: iam.cnrm.cloud.google.com/v1beta1
kind: IAMPolicy
metadata:
  name: storage-admin-policy
spec:
  resourceRef:
    apiVersion: iam.cnrm.cloud.google.com/v1beta1
    kind: IAMServiceAccount
    name: storage-admin
    namespace: cnrm-system
  bindings:
    - role: roles/iam.workloadIdentityUser
      members:
        - serviceAccount:mysten-sui.svc.id.goog[monitoring/test-compactor]
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
apiVersion: generators.kustomize.devopshive.net/v1alpha1
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
  - name: my-iam-policy
    input:
      method: overlay
      files:
        # The same as in the KRM container above, omit it if Exec KRM is used.
        root: /mnt
        sources:
        - input/iam-policy-src.yaml
        destination: input/iam-policy-dst.yaml
    merge:
      strategy: append
    output:
      format: raw
```

## Build

```shell
kustomize build  --enable-alpha-plugins --as-current-user .
```

## Output

```yaml
apiVersion: iam.cnrm.cloud.google.com/v1beta1
kind: IAMPolicy
metadata:
  name: storage-admin-policy
  namespace: cnrm-system
spec:
  bindings:
  - members:
    - serviceAccount:mysten-sui.svc.id.goog[monitoring/thanos-compactor]
    role: roles/iam.workloadIdentityUser
  - members:
    - serviceAccount:mysten-sui.svc.id.goog[monitoring/test-compactor]
    role: roles/iam.workloadIdentityUser
  resourceRef:
    apiVersion: iam.cnrm.cloud.google.com/v1beta1
    kind: IAMServiceAccount
    name: storage-admin
    namespace: cnrm-system
```