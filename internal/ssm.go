package hssm

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// GetSSMParameter gets a parameter from the AWS Simple Systems Manager service.
func GetSSMParameter(svc ssmiface.SSMAPI, name string, defaultValue *string, decrypt bool) (*string, error) {
	regex := "([a-zA-Z0-9\\.\\-_/]*)"
	r, _ := regexp.Compile(regex)

	if match := r.FindString(name); match == "" {
		return nil,
			fmt.Errorf(
				"there is an invalid character in the name of the parameter: %s. It should match %s",
				name,
				regex,
			)
	}
	// Create the request to SSM
	getParameterInput := &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: &decrypt,
	}

	// Get the parameter from SSM
	param, err := svc.GetParameter(getParameterInput)
	// Cast err to awserr.Error to handle specific error codes.
	// TODO(dmizelle): convert this to errors.As, needs aws sdk v2
	aerr, ok := err.(awserr.Error) // nolint:errorlint
	if ok && aerr.Code() == ssm.ErrCodeParameterNotFound && defaultValue != nil {
		return defaultValue, nil
	}

	if err != nil {
		return nil, fmt.Errorf("unable to get SSM parameter: %w", err)
	}

	return param.Parameter.Value, nil
}
