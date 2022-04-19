package config

import (
	"errors"
	"github.com/healthimation/go-aws-config/src/provider"
	"github.com/joho/godotenv"
	"os"
)

var (
	ErrEnvNotSet = errors.New("env not set")
)

const (
	configKeyEnvironment = "HMD_ENVIRONMENT"
	configProviderParam  = "CONFIG_PROVIDER"
	configFilePath       = "CFG"
	configAwsRegion      = "AWS_REGION"
)

func NewConfigProviderFromEnv(defaultServiceName string) (provider.Provider, error) {
	// pull environment from env vars
	cfgFile := os.Getenv(configFilePath)
	if len(cfgFile) > 0 {
		_ = godotenv.Load(cfgFile)
	}
	// pull environment from env vars
	env := os.Getenv(configKeyEnvironment)
	if len(env) == 0 {
		return nil, ErrEnvNotSet
	}
	cfgProvider := os.Getenv(configProviderParam)
	if len(cfgProvider) == 0 {
		cfgProvider = "secrets_manager"
	}
	awsRegion := os.Getenv(configAwsRegion)

	return NewConfigProvider(cfgProvider, &provider.Config{
		Env:         env,
		ServiceName: defaultServiceName,
		AWSRegion:   awsRegion,
	})
}
