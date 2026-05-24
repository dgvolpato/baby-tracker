package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"baby-tracker/internal/auth"

	"github.com/aws/aws-lambda-go/events"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	var body loginRequest
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return respond(http.StatusBadRequest, `{"error":"invalid request"}`, headers), nil
	}

	userID, ok := auth.CheckCredentials(body.Username, body.Password)
	if !ok {
		return respond(http.StatusUnauthorized, `{"error":"invalid credentials"}`, headers), nil
	}

	token, err := auth.SignToken(userID)
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"could not sign token"}`, headers), nil
	}

	out, _ := json.Marshal(map[string]string{"token": token})
	return respond(http.StatusOK, string(out), headers), nil
}
