package main

import (
	"context"
	"net/http"
	"strings"

	"baby-tracker/internal/auth"
	"baby-tracker/internal/handlers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	method := req.RequestContext.HTTP.Method
	path := req.RawPath

	headers := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Authorization, Content-Type",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
	}

	if method == "OPTIONS" {
		return events.APIGatewayV2HTTPResponse{StatusCode: 200, Headers: headers}, nil
	}

	if method == "POST" && path == "/login" {
		return handlers.Login(ctx, req, headers)
	}

	tokenStr := strings.TrimPrefix(req.Headers["authorization"], "Bearer ")
	userID, err := auth.VerifyToken(tokenStr)
	if err != nil || userID == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Headers:    headers,
			Body:       `{"error":"unauthorized"}`,
		}, nil
	}

	switch {
	case method == "GET" && path == "/feedings":
		return handlers.ListFeedings(ctx, req, headers)
	case method == "POST" && path == "/feedings":
		return handlers.CreateFeeding(ctx, req, headers)
	case method == "PUT" && strings.HasPrefix(path, "/feedings/"):
		return handlers.UpdateFeeding(ctx, req, headers)
	case method == "DELETE" && strings.HasPrefix(path, "/feedings/"):
		return handlers.DeleteFeeding(ctx, req, headers)
	default:
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusNotFound,
			Headers:    headers,
			Body:       `{"error":"not found"}`,
		}, nil
	}
}
