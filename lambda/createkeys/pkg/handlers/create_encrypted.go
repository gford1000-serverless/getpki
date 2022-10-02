package handlers

import (
	"createkeys/pkg/util"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gford1000-go/pkigen"
)

var errInvalidArgument = errors.New("unexpected argument provided")

// CreateEncrypted expects to receive a Base64 encoded RSA public key, which it
// will then be used to encrypt the randomly generated RSA key pair of the specified size
func CreateEncrypted(event interface{}, size int) (*events.APIGatewayProxyResponse, error) {
	encData, ok := event.(*keyRequestEvent)
	if !ok {
		return util.NewErrorAPIResponse(500, errInvalidArgument), nil
	}

	k, err := pkigen.UnmarshalPublicKey(
		&pkigen.Base64EncodedRSAKey{
			PublicKey: encData.PublicKey,
		})
	if err != nil {
		return util.NewErrorAPIResponse(400, err), nil
	}

	e, err := pkigen.CreateEncryptedRSAKey(k, size)
	if err != nil {
		return util.NewErrorAPIResponse(400, err), nil
	}

	return util.NewAPIResponse(200, e)
}
