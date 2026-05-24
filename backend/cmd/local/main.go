// Local development server — wraps the Lambda handler in a plain net/http server.
// Run: go run ./cmd/local from the backend directory (after copying .env.example to .env).
package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"baby-tracker/internal/auth"
	"baby-tracker/internal/handlers"

	"github.com/aws/aws-lambda-go/events"
)

func main() {
	loadEnv(".env")

	mux := http.NewServeMux()
	mux.HandleFunc("/", adapt)

	addr := ":3001"
	log.Printf("local server listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, corsMiddleware(mux)))
}

func adapt(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)

	flatHeaders := make(map[string]string, len(r.Header))
	for k, v := range r.Header {
		flatHeaders[strings.ToLower(k)] = v[0]
	}

	req := events.APIGatewayV2HTTPRequest{
		RawPath: r.URL.Path,
		Headers: flatHeaders,
		Body:    string(bodyBytes),
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: r.Method,
			},
		},
	}

	resp, err := lambdaHandler(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range resp.Headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(resp.StatusCode)
	if resp.Body != "" {
		w.Write([]byte(resp.Body)) //nolint:errcheck
	}
}

func lambdaHandler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// loadEnv reads a simple KEY=VALUE file and sets env vars (does not override existing ones).
func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if os.Getenv(strings.TrimSpace(k)) == "" {
			os.Setenv(strings.TrimSpace(k), strings.TrimSpace(v))
		}
	}
}
