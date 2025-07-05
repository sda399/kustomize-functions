# Password generator

This transformer aims to set randomly generated values to `Secrets`.
Currently the follwoing data types can begenerated:
* password
* uuid
* ssh key pairs

## Example
```yaml
# kustomization.yaml
transformers:
  - ./transformer.yaml
secretGenerator:
  - name: uuid-random
    literals:
      - USER_ID=""
      - USER_SECRET=""
    options:
      annotations:
        local.config.zprd/passwords: |-
          - name: USER_SECRET
            type: random
          - name: USER_ID
            type: uuid
---
# transformer.yaml
apiVersion: v1
kind: RandomReplacer
metadata:
  name: randomReplacer
  annotations:
    config.kubernetes.io/function: |-
      container:
        image: ghcr.io/zprd/kustomize-functions/password-generator:v1.0.0
```
