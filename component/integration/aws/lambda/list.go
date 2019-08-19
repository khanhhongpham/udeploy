package lambda

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/turnerlabs/udeploy/component/app"
)

// ListDefinitions ...
func ListDefinitions(instance app.Instance) (map[string]app.Definition, error) {

	svc := lambda.New(session.New())
	o, err := svc.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
		FunctionName: aws.String(instance.FunctionName),
	})
	if err != nil {
		return nil, err
	}

	versions := map[string]app.Definition{}
	for _, funcVersion := range o.Versions {

		revision, err := strconv.ParseInt(*funcVersion.Version, 10, 64)
		if err != nil {
			continue
		}

		version, build, err := extractVersion(instance, funcVersion)
		if err != nil {
			continue
		}

		env := map[string]string{}
		for k, v := range funcVersion.Environment.Variables {
			value := *v
			env[k] = value
		}

		def := app.Definition{
			Description: version,

			Version:  version,
			Build:    build,
			Revision: revision,

			Environment: env,
			Secrets:     map[string]string{},
		}

		versions[def.FormatVersion()] = def
	}

	return versions, nil
}