package util

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func addHeaders() map[string]string {
	return map[string]string{"Content-Type": "application/json"}
}

// NewErrorAPIResponse returns a standard body for errors
func NewErrorAPIResponse(status int, e error) *events.APIGatewayProxyResponse {
	type errorMsg struct {
		Msg string `json:"error"`
	}

	stringBody, _ := json.Marshal(
		&errorMsg{
			Msg: e.Error(),
		})

	return &events.APIGatewayProxyResponse{
		Headers:    addHeaders(),
		StatusCode: status,
		Body:       string(stringBody),
	}
}

// NewAPIResponse formats the response for the API Gateway
func NewAPIResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	stringBody, err := json.Marshal(body)
	if err != nil {
		return NewErrorAPIResponse(500, err), nil
	}

	return &events.APIGatewayProxyResponse{
		Headers:    addHeaders(),
		StatusCode: status,
		Body:       string(stringBody),
	}, nil
}
