package main

import (
	"fmt"
	"os"
	"path/filepath"

	"config.zprd/replacementextra/filters"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
)

func main() {
	functionConfig := &filters.Filter{}
	p := framework.SimpleProcessor{Config: functionConfig, Filter: functionConfig}
	//if err := framework.Execute(p, nil); err != nil {
	//	fmt.Fprintf(os.Stderr, err.Error())
	//	os.Exit(1)
	//}
	cmd := command.Build(p, command.StandaloneDisabled, false)
	addGenerateDockerfile(cmd)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func addGenerateDockerfile(cmd *cobra.Command) {
	gen := &cobra.Command{
		Use:  "gen [DIR]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.WriteFile(filepath.Join(args[0], "Dockerfile"), []byte(`
FROM public.ecr.aws/docker/library/golang:1.24-bullseye AS builder
ENV CGO_ENABLED=0
WORKDIR /go/src/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags '-w -s' -v -o /usr/local/bin/function ./

FROM gcr.io/distroless/static-debian12:latest
COPY --from=builder /usr/local/bin/function /usr/local/bin/function
ENTRYPOINT ["function"]
`), 0600); err != nil {
				return fmt.Errorf("%w", err)
			}
			return nil
		},
	}
	cmd.AddCommand(gen)
}
