package config_test

import (
	"bytes"
	"testing"

	"github.com/nil-nil/ticket/internal/services/config"
	"github.com/stretchr/testify/assert"
)

func TestYamlConfig(t *testing.T) {
	structConfig := config.Config{
		HTTP: struct {
			Port          int    `yaml:"port"`
			ListenAddress string `yaml:"listenAddress"`
		}{
			Port:          8080,
			ListenAddress: "localhost",
		},
		Auth: struct {
			JWT *struct {
				SigningMethod string `yaml:"signingMethod"`
				TokenLifetime uint64 `yaml:"tokenLifetime"`
				PublicKey     string `yaml:"publicKey"`
				PrivateKey    string `yaml:"privateKey"`
			} `yaml:"jwt"`
		}{
			JWT: &struct {
				SigningMethod string `yaml:"signingMethod"`
				TokenLifetime uint64 `yaml:"tokenLifetime"`
				PublicKey     string `yaml:"publicKey"`
				PrivateKey    string `yaml:"privateKey"`
			}{
				SigningMethod: "RS512",
				TokenLifetime: 518500,
				PublicKey:     "testPublicKey\n12345\n",
				PrivateKey:    "testPrivateKey\n67890\n",
			},
		},
	}

	yamlConfig := `
httpServer:
  port: 8080
  listenAddress: localhost
auth:
  jwt:
    signingMethod: RS512
    tokenLifetime: 518500
    publicKey: |
      testPublicKey
      12345
    privateKey: |
      testPrivateKey
      67890
`

	b := bytes.NewBufferString(yamlConfig)

	config, err := config.GetConfig(b)

	assert.NoError(t, err)
	assert.Equal(t, structConfig, config)
}
