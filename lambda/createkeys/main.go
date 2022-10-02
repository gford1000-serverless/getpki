package main

import (
	"context"
	"createkeys/pkg/handlers"
	"createkeys/pkg/util"
	"errors"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Default to 2048
var size = 2048

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
}

var errInvalidHTTPMethod = errors.New("only GET or POST requests supported")

func handleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	// Always respond to a GET request
	if event.HTTPMethod == "GET" {
		return handlers.CreateUnencrypted(size)
	}

	// Otherwise must be a POST request, containing a valid body
	if event.HTTPMethod != "POST" {
		return util.NewErrorAPIResponse(400, errInvalidHTTPMethod), nil
	}

	body, err := handlers.Unpack(event.Body)
	if err != nil {
		return util.NewErrorAPIResponse(400, err), nil
	}

	return handlers.CreateEncrypted(body, size)
}

func main() {
	lambda.Start(handleRequest)
}
