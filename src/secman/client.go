package secman

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"time"
)

var (
	ErrCodeDecryptionFailure         = errors.New("secrets manager can't decrypt the protected secret text using the provided KMS key")
	ErrCodeInternalServiceError      = errors.New("an error occurred on the server side")
	ErrCodeInvalidParameterException = errors.New("you provided an invalid value for a parameter")
	ErrCodeInvalidRequestException   = errors.New("you provided a parameter value that is not valid for the current state of the resource")
	ErrCodeResourceNotFoundException = errors.New("we can't find the resource that you asked for")
	ErrUnknown                       = errors.New("unknown error")
)

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

type secman struct {
	man          *secretsmanager.SecretsManager
	env          string
	versionStage string
}

func NewConfigProvider(sess *session.Session, env, region string) (Provider, error) {
	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(region))
	return &secman{
		man:          svc,
		env:          env,
		versionStage: "AWSCURRENT",
	}, nil
}

func (svc *secman) Import(data []byte) error {
	return nil
}
func (svc *secman) Initialize() error {
	return nil
}
func (svc *secman) Get(key string) ([]byte, error) {
	return nil, nil
}
func (svc *secman) Put(key string, value []byte) error {
	return nil
}

func (svc *secman) MustGetString(key string) string {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String(svc.versionStage), // VersionStage defaults to AWSCURRENT if unspecified
	}
	result, err := svc.man.GetSecretValue(input)
	if err = handleAwsErr(err); err != nil {
		panic(err)
	}
	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
		return secretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			panic(err)
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		return decodedBinarySecret
	}
}

func (svc *secman) MustGetBool(key string) bool {
	panic(ErrUnknown)
}

func (svc *secman) MustGetInt(key string) int {
	panic(ErrUnknown)
}
func (svc *secman) MustGetDuration(key string) time.Duration {
	panic(ErrUnknown)
}

func handleAwsErr(err error) error {
	if err == nil {
		return nil
	}
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case secretsmanager.ErrCodeDecryptionFailure:
			// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
			fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			return ErrCodeDecryptionFailure

		case secretsmanager.ErrCodeInternalServiceError:
			// An error occurred on the server side.
			fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			return ErrCodeInternalServiceError

		case secretsmanager.ErrCodeInvalidParameterException:
			// You provided an invalid value for a parameter.
			fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			return ErrCodeInvalidParameterException

		case secretsmanager.ErrCodeInvalidRequestException:
			// You provided a parameter value that is not valid for the current state of the resource.
			fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			return ErrCodeInvalidRequestException

		case secretsmanager.ErrCodeResourceNotFoundException:
			// We can't find the resource that you asked for.
			fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			return ErrCodeResourceNotFoundException

		default:
			fmt.Println(aerr.Error())
			return ErrUnknown
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return ErrUnknown
	}
}
