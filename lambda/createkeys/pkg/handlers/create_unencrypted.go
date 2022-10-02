package handlers

import (
	"createkeys/pkg/util"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gford1000-go/pkigen"
)

// CreateUnencrypted generates a new pair of RSA keys, of the specified size
func CreateUnencrypted(size int) (*events.APIGatewayProxyResponse, error) {
	b, err := pkigen.CreateEncodedRSAKey(size)
	if err != nil {
		return util.NewErrorAPIResponse(500, err), nil
	}

	return util.NewAPIResponse(200, b)
}
