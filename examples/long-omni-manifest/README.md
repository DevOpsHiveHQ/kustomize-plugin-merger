<!-- omit in toc -->
# Example - Organize long manifests into smaller ones

- [Use case](#use-case)
- [Input](#input)
- [Manifest](#manifest)
- [Build](#build)
- [Output](#output)

## Use case

In some use cases, you could have a really long YAML manifest, and it's hard to read. You can split that file
and use the Merger `patch` input method to make it a single manifest again.

An example of that is working with [Crossplane Compositions](https://docs.crossplane.io/latest/concepts/compositions/).
A single composition file could be hundreds of lines. Let's take an example:
the [EKS composition](https://github.com/upbound/platform-ref-aws/tree/main/package/cluster/eks)
at [AWS Reference Platform](https://github.com/upbound/platform-ref-aws).

The EKS composition is about 300 lines with 14 resources (in some other scenarios, it could be even more!),
so Kustomize and Merger could be used together to orchestrate producing that file.
We will use the 1st 2 out of the 14 resources for demonstration since it's the same logic anyway.


## Input

The [package/cluster/eks/composition.yaml](https://github.com/upbound/platform-ref-aws/blob/c972e0cd2fe91ffdb0b181250832fae74e473b10/package/cluster/eks/composition.yaml).


```yaml
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: xeks.aws.platformref.upbound.io
  labels:
    provider: aws
spec:
  writeConnectionSecretsToNamespace: upbound-system
  compositeTypeRef:
    apiVersion: aws.platformref.upbound.io/v1alpha1
    kind: XEKS
  resources:
    - base:
        apiVersion: iam.aws.upbound.io/v1beta1
        kind: Role
        metadata:
          labels:
            role: controlplane
        spec:
          forProvider:
            assumeRolePolicy: |
              {
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Effect": "Allow",
                        "Principal": {
                            "Service": [
                                "eks.amazonaws.com"
                            ]
                        },
                        "Action": [
                            "sts:AssumeRole"
                        ]
                    }
                ]
              }
      name: controlplaneRole
    - base:
        apiVersion: iam.aws.upbound.io/v1beta1
        kind: RolePolicyAttachment
        spec:
          forProvider:
            policyArn: arn:aws:iam::aws:policy/AmazonEKSClusterPolicy
            roleSelector:
              matchControllerRef: true
              matchLabels:
                role: controlplane
      name: clusterRolePolicyAttachment
    ##############################
    # More 12 resources redacted #
    ##############################
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
  - name: my-eks-composition
    input:
      # Here patch is used to produce 1 resource from all sources and destination.
      method: patch
      files:
        # The same as in the KRM container above, omit it if Exec KRM is used.
        root: /mnt
        sources:
        # Just the 1st 2 items in the composition resources as an example.
        - input/controlplaneRole.yaml
        - input/clusterRolePolicyAttachment.yaml
        destination: input/composition-base.yaml
    merge:
      # Add all list items to the same path.
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
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  labels:
    provider: aws
  name: xeks.aws.platformref.upbound.io
spec:
  compositeTypeRef:
    apiVersion: aws.platformref.upbound.io/v1alpha1
    kind: XEKS
  resources:
  - base:
      apiVersion: iam.aws.upbound.io/v1beta1
      kind: Role
      metadata:
        labels:
          role: controlplane
      spec:
        forProvider:
          assumeRolePolicy: |
            {
              "Version": "2012-10-17",
              "Statement": [
                  {
                      "Effect": "Allow",
                      "Principal": {
                          "Service": [
                              "eks.amazonaws.com"
                          ]
                      },
                      "Action": [
                          "sts:AssumeRole"
                      ]
                  }
              ]
            }
    name: controlplaneRole
  - base:
      apiVersion: iam.aws.upbound.io/v1beta1
      kind: RolePolicyAttachment
      spec:
        forProvider:
          policyArn: arn:aws:iam::aws:policy/AmazonEKSClusterPolicy
          roleSelector:
            matchControllerRef: true
            matchLabels:
              role: controlplane
    name: clusterRolePolicyAttachment
  writeConnectionSecretsToNamespace: upbound-system
```