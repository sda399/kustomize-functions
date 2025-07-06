# Replacement extra

This Function is a copy of the Kustomize `replacement` builtin transformer
with the added support for regex expressions to select `targets`

This makes it possible to select multiple resources within the same target block

## Example
Add service namespace to an environment variable in Deployments, DaemonSet, StatefulSet, declared in kustomize resources

```yaml
apiVersion: kustomize-functions.zprd.io/v1
kind: ReplacementExtra

metadata:
  name: replacementServiceUrl
  annotations:
    config.kubernetes.io/function: |-
      container:
        image: ghcr.io/sda399/kustomize-functions/replacement-extra:v1.0.0

replacements:
  - source:
      kind: Service
      name: svc1
      fieldPath: metadata.namespace
    targets:
      - select:
          kind: Deployment|.*Set
        fieldPaths:
          - spec.template.spec.containers.*.env.[name=APP1_SERVICE].value
        options:
          index: 99
          delimiter: "."
```
