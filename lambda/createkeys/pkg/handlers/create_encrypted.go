package handlers

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gford1000-go/pkigen"
	util "github.com/gford1000-serverless/util/events"
)

var errInvalidArgument = errors.New("unexpected argument provided")

// CreateEncrypted expects to receive a Base64 encoded RSA public key, which it
// will then be used to encrypt the randomly generated RSA key pair of the specified size
func CreateEncrypted(event interface{}, size int, responder *util.GatewayProxyResponder) (*events.APIGatewayProxyResponse, error) {
	encData, ok := event.(*keyRequestEvent)
	if !ok {
		return responder.NewErrorAPIResponse(500, errInvalidArgument), nil
	}

	k, err := pkigen.UnmarshalPublicKey(
		&pkigen.Base64EncodedRSAKey{
			PublicKey: encData.PublicKey,
		})
	if err != nil {
		return responder.NewErrorAPIResponse(400, err), nil
	}

	e, err := pkigen.CreateEncryptedRSAKey(k, size)
	if err != nil {
		return responder.NewErrorAPIResponse(400, err), nil
	}

	return responder.NewAPIResponse(200, e)
}
