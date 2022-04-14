package secman

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/healthimation/go-aws-config/src/provider"
	"strconv"
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

type secman struct {
	man          *secretsmanager.SecretsManager
	env          string
	versionStage string
}

func NewConfigProvider(sess *session.Session, env, region string) provider.Provider {
	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(region))
	return &secman{
		man:          svc,
		env:          env,
		versionStage: "AWSCURRENT",
	}
}

func (svc *secman) Import(data []byte) error {
	return nil
}

func (svc *secman) Initialize() error {
	return nil
}

func (svc *secman) Get(key string) ([]byte, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String(svc.versionStage), // VersionStage defaults to AWSCURRENT if unspecified
	}
	result, err := svc.man.GetSecretValue(input)
	if err = handleAwsErr(err); err != nil {
		return nil, err
	}
	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
		return []byte(secretString), nil
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		_, err = base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			return nil, err
		}
		return decodedBinarySecretBytes, nil
	}
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
		l, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			panic(err)
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:l])
		return decodedBinarySecret
	}
}

func (svc *secman) MustGetBool(key string) bool {
	val := svc.MustGetString(key)
	ret, err := strconv.ParseBool(val)
	if err != nil {
		panic(err)
	}
	return ret
}

func (svc *secman) MustGetInt(key string) int {
	val := svc.MustGetString(key)
	ret, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return ret
}
func (svc *secman) MustGetDuration(key string) time.Duration {
	val := svc.MustGetString(key)
	ret, err := time.ParseDuration(val)
	if err != nil {
		panic(err)
	}
	return ret
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
