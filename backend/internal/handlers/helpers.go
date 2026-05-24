package handlers

import "github.com/aws/aws-lambda-go/events"

func respond(status int, body string, headers map[string]string) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers:    headers,
		Body:       body,
	}
}
