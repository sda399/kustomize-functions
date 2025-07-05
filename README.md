# KRM functions aka Kustomize functions

This is a collection of KRM functions that can be used in Kustomize projects to extend Kustomize builtin features

## List of available functions

* **[replacement extra](./functions/replacement-extra/README.md)**: the kustomize builtin replacement with extra features: regex support
* **[password generator](./functions/password-generator/README.md)**: inject random data in your secrets: password, ssh key pair, uuid

## Usage example

The following shows how to integrate replacement-extra ina kustomization file.
More examples are available in the Kustomize project repository

### kustomization.yaml
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - resources.yaml

transformers:
  - transformer.yaml
```


### transformer.yaml using a Docker image
Functions can be configured using a docker image or an executable file

```yaml
apiVersion: kustomize-functions.zprd.io/v1
kind: ReplacementExtra

metadata:
  name: replacementServiceUrl
  annotations:
    config.kubernetes.io/function: |-
      container:
        image: ghcr.io/zprd/kustomize-functions/replacement-extra:v1.0.0

replacements:
...
```

### transformer.yaml as local executable file
```sh
curl "$URL/replacement-extra" -o /opt/krm-functions/replacement-extra
export PATH=$PATH:/opt/krm-functions
```

```yaml
apiVersion: kustomize-functions.zprd.io/v1
kind: ReplacementExtra

metadata:
  name: replacementServiceUrl
  annotations:
    config.kubernetes.io/function: |-
      exec:
        path: replacement-extra

replacements:
...
```

#### Build
build when invoking a docker image:

    kustomize build --enable-alpha-plugins ./

alternatively when using a executable file:

    kustomize build --enable-alpha-plugins --enable-exec ./
