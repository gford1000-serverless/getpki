package main

import (
	"context"
	"createkeys/pkg/handlers"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/gford1000-go/pkigen"
	util "github.com/gford1000-serverless/util/events"
)

// Default to 2048
var size = 2048

var responder *util.GatewayProxyResponder

func init() {
	// Size can be changed via a Lambda environment variable
	s := os.Getenv("RSA_SIZE")
	if len(s) > 0 {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		size = i
	}

	// Default to no signing
	responder = util.NewGatewayProxyResponder(nil)
	responder.AddHeader("Content-Type", "application/json")

	// See whether a signing key has been provided
	s = os.Getenv("SIGNING_KEY")
	keyID := os.Getenv("SIGNING_KMS_ID")
	if len(s) > 0 && len(keyID) > 0 {
		// This is a KMS encrypted, JSON serialised Base64EncodedRSAKey
		// that contains a private_key attribute
		sb, err := base64.URLEncoding.DecodeString(s)
		if err != nil {
			fmt.Printf("Unable to decode SIGNING_KEY: %s\n", err)
			return
		}

		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			fmt.Printf("Unable to load default context: %s\n", err)
			return
		}

		client := kms.NewFromConfig(cfg)
		input := &kms.DecryptInput{
			CiphertextBlob: sb,
			KeyId:          &keyID,
		}

		result, err := client.Decrypt(context.TODO(), input)
		if err != nil {
			fmt.Printf("Error decrypting with key (%s): %s\n", keyID, err)
			return
		}

		var b pkigen.Base64EncodedRSAKey
		err = json.Unmarshal(result.Plaintext, &b)
		if err != nil {
			fmt.Printf("Error unmarshaling plaintext JSON: %s\n", err)
			return
		}

		signingPrivateKey, err := pkigen.UnmarshalPrivateKey(&b)
		if err != nil {
			fmt.Printf("Error extracting private key for siging: %s\n", err)
			return
		}

		responder = util.NewGatewayProxyResponder(signingPrivateKey)
		responder.AddHeader("Content-Type", "application/json")
	}
}

var errInvalidHTTPMethod = errors.New("only GET or POST requests supported")

func handleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	// Always respond to a GET request
	if event.HTTPMethod == "GET" {
		return handlers.CreateUnencrypted(size, responder)
	}

	// Otherwise must be a POST request, containing a valid body
	if event.HTTPMethod != "POST" {
		return responder.NewErrorAPIResponse(400, errInvalidHTTPMethod), nil
	}

	body, err := handlers.Unpack(event.Body)
	if err != nil {
		return responder.NewErrorAPIResponse(400, err), nil
	}

	return handlers.CreateEncrypted(body, size, responder)
}

func main() {
	lambda.Start(handleRequest)
}
