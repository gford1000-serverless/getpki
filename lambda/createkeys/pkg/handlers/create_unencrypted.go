package handlers

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/gford1000-go/pkigen"
	util "github.com/gford1000-serverless/util/events"
)

// CreateUnencrypted generates a new pair of RSA keys, of the specified size
func CreateUnencrypted(size int, responder *util.GatewayProxyResponder) (*events.APIGatewayProxyResponse, error) {
	b, err := pkigen.CreateEncodedRSAKey(size)
	if err != nil {
		return responder.NewErrorAPIResponse(500, err), nil
	}

	return responder.NewAPIResponse(200, b)
}
