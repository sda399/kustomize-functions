package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"config.zprd/replacementextra/filters"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/kio"
)

func Test(t *testing.T) {
	functionConfig := &filters.Filter{}
	p := framework.SimpleProcessor{Config: functionConfig, Filter: functionConfig}
	if err := framework.Execute(p, nil); err != nil {
		os.Exit(1)
	}
	out := new(bytes.Buffer)
	rw := &kio.ByteReadWriter{Reader: bytes.NewBufferString(`
apiVersion: config.kubernetes.io/v1
kind: ResourceList
items:
- apiVersion: v1
  kind: Deployment
  metadata:
    name: deploy
  spec:
    template:
      spec:
        containers:
        - image: nginx:1.7.9
          name: nginx-tagged
        - image: postgres:1.8.0
          name: postgresdb
functionConfig:
  replacements:
  - source:
      kind: Deployment
      name: deploy
      fieldPath: spec.template.spec.containers.0.image
    targets:
    - select:
        kind: Deployment
        name: deploy
      fieldPaths:
      - spec.template.spec.containers.1.image
`),
		Writer: out}

	require.NoError(t, framework.Execute(p, rw))
	require.Equal(t, strings.TrimSpace(`
apiVersion: config.kubernetes.io/v1
kind: ResourceList
items:
- apiVersion: v1
  kind: Deployment
  metadata:
    name: deploy
  spec:
    template:
      spec:
        containers:
        - image: nginx:1.7.9
          name: nginx-tagged
        - image: nginx:1.7.9
          name: postgresdb
functionConfig:
  replacements:
  - source:
      kind: Deployment
      name: deploy
      fieldPath: spec.template.spec.containers.0.image
    targets:
    - select:
        kind: Deployment
        name: deploy
      fieldPaths:
      - spec.template.spec.containers.1.image
`), strings.TrimSpace(out.String()))
}
