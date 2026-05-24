// Local development server — wraps the Lambda handler in a plain net/http server.
// Run: go run ./cmd/local (after copying .env.example to .env)
package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"baby-tracker/internal/router"

	"github.com/aws/aws-lambda-go/events"
)

func main() {
	loadEnv(".env")

	mux := http.NewServeMux()
	mux.HandleFunc("/", adapt)

	addr := ":3001"
	log.Printf("local server listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
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

	resp, err := router.Handle(r.Context(), req)
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

// keep corsMiddleware to avoid unused warning — router.Handle handles CORS internally
var _ = corsMiddleware
