package chainsuite

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const prefix = "TEST"

type Environment struct {
	DockerRegistry        string `envconfig:"DOCKER_REGISTRY" default:"ghcr.io/mantra-chain"`
	MantraImageName       string `envconfig:"MANTRA_IMAGE_NAME" default:"mantrachain"`
	OldMantraImageVersion string `envconfig:"OLD_MANTRA_IMAGE_VERSION"`
	NewMantraImageVersion string `envconfig:"NEW_MANTRA_IMAGE_VERSION"`
	UpgradeName           string `envconfig:"UPGRADE_NAME"`
}

func GetEnvironment() Environment {
	var env Environment
	if err := envconfig.Process(prefix, &env); err != nil {
		panic(fmt.Errorf("failed to process environment variables: %w", err))
	}
	return env
}
