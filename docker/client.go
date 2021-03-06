package docker

import (
	"os"

	"github.com/docker/docker/client"
)

// Client docker client
var Client *Docker

// Docker docker client
type Docker struct {
	*client.Client
}

// ClientConfig docker client config
type ClientConfig struct {
	endpoint      string
	certPath      string
	keyPath       string
	caPath        string
	engineVersion string
}

// NewClientConfig create docker client config
func NewClientConfig(endpoint, cert, key, ca, engineVersion string) *ClientConfig {
	return &ClientConfig{
		endpoint:   endpoint,
		certPath:   cert,
		keyPath:    key,
		caPath:     ca,
		engineVersion: engineVersion,
	}
}

// NewDocker create new docker client
func NewDocker(config *ClientConfig) *Docker {
	if os.Getenv("DOCKER_HOST") != "" {
		client, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion(config.engineVersion))
		if err != nil {
			panic(err)
		}

		return &Docker{client}
	}

	if config.caPath != "" &&
		config.certPath != "" &&
		config.keyPath != "" {
		client, err := client.NewClientWithOpts(client.WithTLSClientConfig(config.caPath, config.certPath, config.keyPath),
			client.WithHost(config.endpoint),
			client.WithVersion(config.engineVersion))

		if err != nil {
			panic(err)
		}

		return &Docker{client}
	}

	client, err := client.NewClientWithOpts(client.WithHost(config.endpoint), client.WithVersion(config.engineVersion))
	if err != nil {
		panic(err)
	}

	Client = &Docker{client}

	return Client
}
