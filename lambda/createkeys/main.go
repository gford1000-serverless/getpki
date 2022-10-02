package main

import (
	"context"
	"createkeys/pkg/handlers"
	"errors"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

	// Load and set signing key here, to sign responses
	responder = util.NewGatewayProxyResponder(nil)
	responder.AddHeader("Content-Type", "application/json")
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
