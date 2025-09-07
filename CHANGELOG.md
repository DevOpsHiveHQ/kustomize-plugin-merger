# Changelog

## [2.2.0](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/compare/v2.1.1...v2.2.0) (2025-09-07)


### Features

* support annotations ([#91](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/issues/91)) ([88b5b0f](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/commit/88b5b0f1343097e928c789b20301baecbbab6a4c))


### Bug Fixes

* **deps:** update all non-major dependencies ([#89](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/issues/89)) ([d4e1029](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/commit/d4e1029237c945e031c1d359a2100c1f1f660350))

## [2.1.1](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/compare/v2.1.0...v2.1.1) (2024-09-14)


### Refactors

* run merger with non-root user by default ([2faee16](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/commit/2faee16f4d9337738bc60b508f4cbbcede942e55))


### Continuous Integration

* publish temporary docker image ([cf83a24](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/commit/cf83a248ad646611cfca79635c0990c3667d385b))
* update golangci-lint-action version ([f25736e](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/commit/f25736ed665e0e6e7b7898200f005e5123e4c6a9))

## [2.1.0](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/compare/v2.0.0...v2.1.0) (2024-08-20)


### Features

* support docker image arm64 arch ([#79](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/issues/79)) ([f23287d](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/commit/f23287d9614e273fd006a8af89e9e8ea268dd56e))

## [2.0.0](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/compare/v1.1.0...v2.0.0) (2024-02-11)


### ⚠ BREAKING CHANGES

* rename apiVersion group name ([#75](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/issues/75))

### refactor

* rename apiVersion group name ([#75](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/issues/75)) ([d61661b](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/commit/d61661b12111a7f422535f20815cbee61aab4276))


### Bug Fixes

* **deps:** update all non-major dependencies ([c040f6d](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/commit/c040f6d0d68b188006d14ea0b742868fdea6b815))
* **deps:** update k8s.io/kube-openapi digest to a781456 ([c7a678c](https://github.com/DevOpsHiveHQ/kustomize-plugin-merger/commit/c7a678c88af453270d0dfd063610294b4a8b0402))

## [1.1.0](https://github.com/aabouzaid/kustomize-plugin-merger/compare/v1.0.0...v1.1.0) (2023-09-17)


### Features

* support ConfigMap and Secret as an output ([#34](https://github.com/aabouzaid/kustomize-plugin-merger/issues/34)) ([62e552e](https://github.com/aabouzaid/kustomize-plugin-merger/commit/62e552e49738bedf7dcaf32d49e3ea6ef70caf5b))


### Bug Fixes

* **deps:** update all non-major dependencies ([3f5c70c](https://github.com/aabouzaid/kustomize-plugin-merger/commit/3f5c70c33f2a89927f52292000cb18e3f280d790))


### Continuous Integration

* unify build and release packages ([#26](https://github.com/aabouzaid/kustomize-plugin-merger/issues/26)) ([8381945](https://github.com/aabouzaid/kustomize-plugin-merger/commit/8381945bd297b80431288d83a24ad3398404b6be))

## 1.0.0 (2023-09-10)


### Features

* first working version of merger ([bd0ef52](https://github.com/aabouzaid/kustomize-plugin-merger/commit/bd0ef52ef46381c13ab8ebcf1cc6970328910b64))


### Documentation

* update the main readme and use cases ([1d2e9d4](https://github.com/aabouzaid/kustomize-plugin-merger/commit/1d2e9d45d3aa8ca3bae49ef11361cf4ed43338d9))
* use case - append CustomResources lists without schema ([4be2cc1](https://github.com/aabouzaid/kustomize-plugin-merger/commit/4be2cc168219fe040a556a02cfdd8d76b9c9695e))
* use case - Generate multiple manifests from a single file ([37b9ce1](https://github.com/aabouzaid/kustomize-plugin-merger/commit/37b9ce1c7a8d4daabc5a6762d1d5e3a1566c1e02))
* use case - Organize long manifests into smaller ones ([dad79e0](https://github.com/aabouzaid/kustomize-plugin-merger/commit/dad79e0a9e19f46d261c6ed288d25756d59ab559))


### Continuous Integration

* add go release workflow ([80f6ff4](https://github.com/aabouzaid/kustomize-plugin-merger/commit/80f6ff4b0cf678e902b19eebd0a85db15ed88bde))
* go linter config and ci workflow ([#1](https://github.com/aabouzaid/kustomize-plugin-merger/issues/1)) ([8a3c75e](https://github.com/aabouzaid/kustomize-plugin-merger/commit/8a3c75e1d9eccc9f5ef1211f033b2a5e5d2338cf))
* go releaser ([#2](https://github.com/aabouzaid/kustomize-plugin-merger/issues/2)) ([df20b82](https://github.com/aabouzaid/kustomize-plugin-merger/commit/df20b8281c8107832f0f66d0f857f5f32d558215))
* renovate for dependency update ([#4](https://github.com/aabouzaid/kustomize-plugin-merger/issues/4)) ([61ea189](https://github.com/aabouzaid/kustomize-plugin-merger/commit/61ea189dca345c24e4141bc479da6eb1becfd3aa))
* **security:** apply best practices by step security ([#8](https://github.com/aabouzaid/kustomize-plugin-merger/issues/8)) ([fe208e0](https://github.com/aabouzaid/kustomize-plugin-merger/commit/fe208e06b1d35c454cddcec881d6c282b1691a0a))
* **security:** codeql analysis workflow ([#3](https://github.com/aabouzaid/kustomize-plugin-merger/issues/3)) ([246696c](https://github.com/aabouzaid/kustomize-plugin-merger/commit/246696c73b5c0635ca421727d6e88035a02eb80a))
* **security:** Scorecard security workflow ([52c4496](https://github.com/aabouzaid/kustomize-plugin-merger/commit/52c4496f409c13536c8adf6f3539eb0f1c605944))
