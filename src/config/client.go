package config

import (
	"github.com/healthimation/go-aws-config/src/awsconfig"
	"github.com/healthimation/go-aws-config/src/env"
	"github.com/healthimation/go-aws-config/src/provider"
	"github.com/healthimation/go-aws-config/src/secman"
	"time"
)

type configProvider struct {
	prvd provider.Provider
}

func NewConfigProvider(source string, cfg *provider.Config) (provider.Provider, error) {
	var prvd provider.Provider
	switch source {
	case "parameter_store":
		prvd = awsconfig.NewAWSLoader(cfg.Env, cfg.ServiceName)
	case "secrets_manager":
		prvd = secman.NewConfigProvider(cfg.AWSSession, cfg.Env, cfg.AWSRegion)
	case "env":
		prvd = env.NewConfigProvider()
	}
	return &configProvider{prvd}, nil
}

func (svc *configProvider) Import(data []byte) error {
	return nil
}

func (svc *configProvider) Initialize() error {
	return nil
}

func (svc *configProvider) Get(key string) ([]byte, error) {
	return svc.prvd.Get(key)
}

func (svc *configProvider) Put(key string, value []byte) error {
	return nil
}

func (svc *configProvider) MustGetString(key string) string {
	return svc.prvd.MustGetString(key)
}

func (svc *configProvider) MustGetBool(key string) bool {
	return svc.prvd.MustGetBool(key)
}

func (svc *configProvider) MustGetInt(key string) int {
	return svc.prvd.MustGetInt(key)
}

func (svc *configProvider) MustGetDuration(key string) time.Duration {
	return svc.prvd.MustGetDuration(key)
}
