package provider

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/session"
	"time"
)

var (
	ErrConfigNotFound = errors.New("config not found")
)

type Config struct {
	Env         string
	ServiceName string
	AWSSession  *session.Session
	AWSRegion   string
}

type Provider interface {
	Import(data []byte) error
	Initialize() error
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error

	// Must functions will panic if they can't do what is requested.
	// They are maingly meant for use with configs that are required for an app to start up
	MustGetString(key string) string
	MustGetBool(key string) bool
	MustGetInt(key string) int
	MustGetDuration(key string) time.Duration

	//TODO add array support?
}
