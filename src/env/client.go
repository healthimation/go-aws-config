package env

import (
	"github.com/healthimation/go-aws-config/src/provider"
	"os"
	"strconv"
	"time"
)

type envProvider struct{}

func NewConfigProvider() provider.Provider {
	return &envProvider{}
}

func (svc *envProvider) Import(data []byte) error {
	return nil
}

func (svc *envProvider) Initialize() error {
	return nil
}

func (svc *envProvider) Get(key string) ([]byte, error) {
	d, ok := os.LookupEnv(key)
	if !ok {
		return nil, provider.ErrConfigNotFound
	}
	return []byte(d), nil
}

func (svc *envProvider) Put(key string, value []byte) error {
	return nil
}

func (svc *envProvider) MustGetString(key string) string {
	d, ok := os.LookupEnv(key)
	if !ok {
		panic(provider.ErrConfigNotFound)
	}
	return d
}

func (svc *envProvider) MustGetBool(key string) bool {
	val := svc.MustGetString(key)
	ret, err := strconv.ParseBool(val)
	if err != nil {
		panic(err)
	}
	return ret
}

func (svc *envProvider) MustGetInt(key string) int {
	val := svc.MustGetString(key)
	ret, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return ret
}

func (svc *envProvider) MustGetDuration(key string) time.Duration {
	val := svc.MustGetString(key)
	ret, err := time.ParseDuration(val)
	if err != nil {
		panic(err)
	}
	return ret
}
