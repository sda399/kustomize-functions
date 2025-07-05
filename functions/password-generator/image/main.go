package main

import (
	"crypto"
	"crypto/ed25519"
	uid "github.com/google/uuid"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"

	//"crypto/rsa"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/ssh"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type API struct {
}

type Config struct {
	Type       string `yaml:"type,omitempty"`
	Name       string `yaml:"name,omitempty"`
	PublicKey  string `yaml:"publicKey,omitempty"`
	PrivateKey string `yaml:"privateKey,omitempty"`
}

func main() {
	functionConfig := &API{}
	fn := func(items []*yaml.RNode) ([]*yaml.RNode, error) {
		for i := range items {
			meta, err := items[i].GetMeta()
			if err != nil {
				return nil, err
			}
			if meta.Kind == "Secret" && meta.APIVersion == "v1" {
				const key = "local.config.zprd/passwords"
				if v, ok := meta.Annotations[key]; ok {
					cf, err := readConfigs(v)
					if err != nil {
						return nil, err
					}

					m := items[i].GetDataMap()
					for _, c := range *cf {
						switch c.Type {
						case "ssh":
							pub, priv, err := makeSSHKeyPair()
							if err != nil {
								fmt.Fprintf(os.Stderr, "failed to create ssh key pair: %v", v)
							}
							if c.PublicKey != "" {
								m[c.PublicKey] = base64.StdEncoding.EncodeToString([]byte(pub))
							}
							if c.PrivateKey != "" {
								m[c.PrivateKey] = base64.StdEncoding.EncodeToString([]byte(priv))
							}
							continue
						case "uuid":
							if c.Name != "" {
								m[c.Name] = base64.StdEncoding.EncodeToString([]byte(uid.New().String()))
							}
							continue
						case "random":
						default:
							if c.Name != "" {
								m[c.Name] = base64.StdEncoding.EncodeToString([]byte(random1(32)))
							}
							continue
						}
					}

					items[i].SetDataMap(m)
				}

				delete(meta.Annotations, key)
				if err := items[i].SetAnnotations(meta.Annotations); err != nil {
					return nil, err
				}
			}
		}

		return items, nil
	}

	p := framework.SimpleProcessor{Config: functionConfig, Filter: kio.FilterFunc(fn)}

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

func readConfigs(configs string) (*[]Config, error) {
	c := &[]Config{}
	if err := yaml.Unmarshal([]byte(configs), c); err != nil {
		return nil, fmt.Errorf("in file %s: %w", configs, err)
	}
	return c, nil
}

func random1(length int) string {
	res, err := password.Generate(length, 10, 0, false, true)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	return res
}
func makeSSHKeyPair() (string, string, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	p, err := ssh.MarshalPrivateKey(crypto.PrivateKey(priv), "")
	if err != nil {
		panic(err)
	}
	privateKeyPem := pem.EncodeToMemory(p)
	privateKeyString := string(privateKeyPem)
	publicKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		panic(err)
	}
	publicKeyString := "ssh-ed25519" + " " + base64.StdEncoding.EncodeToString(publicKey.Marshal())
	//fmt.Fprintln(os.Stderr, "Private_Key: %v", privateKeyString)
	//fmt.Fprintln(os.Stderr, "Public_Key: %v", publicKeyString)
	return publicKeyString, privateKeyString, nil
}
