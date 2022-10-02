package handlers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/gford1000-go/pkigen"
	util "github.com/gford1000-serverless/util/events"
)

func AddResponseHeaders() map[string]string {
	return map[string]string{"Content-Type": "application/json"}
}

// CreateUnencrypted generates a new pair of RSA keys, of the specified size
func CreateUnencrypted(size int) (*events.APIGatewayProxyResponse, error) {
	b, err := pkigen.CreateEncodedRSAKey(size)
	if err != nil {
		return util.NewErrorAPIResponse(500, AddResponseHeaders, err), nil
	}

	return util.NewAPIResponse(200, AddResponseHeaders, b)
}
