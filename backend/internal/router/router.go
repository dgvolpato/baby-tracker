// Package router contains the shared request-dispatch logic used by both the
// Lambda handler (cmd/api) and the local dev server (cmd/local).
package router

import (
	"context"
	"net/http"
	"strings"

	"baby-tracker/internal/auth"
	"baby-tracker/internal/handlers"

	"github.com/aws/aws-lambda-go/events"
)

func Handle(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
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
	// Feedings
	case method == "GET" && path == "/feedings":
		return handlers.ListFeedings(ctx, req, headers)
	case method == "POST" && path == "/feedings":
		return handlers.CreateFeeding(ctx, req, headers)
	case method == "PUT" && strings.HasPrefix(path, "/feedings/"):
		return handlers.UpdateFeeding(ctx, req, headers)
	case method == "DELETE" && strings.HasPrefix(path, "/feedings/"):
		return handlers.DeleteFeeding(ctx, req, headers)

	// Diapers
	case method == "GET" && path == "/diapers":
		return handlers.ListDiapers(ctx, req, headers)
	case method == "POST" && path == "/diapers":
		return handlers.CreateDiaper(ctx, req, headers)
	case method == "PUT" && strings.HasPrefix(path, "/diapers/"):
		return handlers.UpdateDiaper(ctx, req, headers)
	case method == "DELETE" && strings.HasPrefix(path, "/diapers/"):
		return handlers.DeleteDiaper(ctx, req, headers)

	// Measurements
	case method == "GET" && path == "/measurements":
		return handlers.ListMeasurements(ctx, req, headers)
	case method == "POST" && path == "/measurements":
		return handlers.CreateMeasurement(ctx, req, headers)
	case method == "PUT" && strings.HasPrefix(path, "/measurements/"):
		return handlers.UpdateMeasurement(ctx, req, headers)
	case method == "DELETE" && strings.HasPrefix(path, "/measurements/"):
		return handlers.DeleteMeasurement(ctx, req, headers)

	default:
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusNotFound,
			Headers:    headers,
			Body:       `{"error":"not found"}`,
		}, nil
	}
}
