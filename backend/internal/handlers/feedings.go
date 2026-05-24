package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"baby-tracker/internal/models"
	"baby-tracker/internal/store"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

func ListFeedings(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	feedings, err := s.ListFeedings(ctx)
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"query failed"}`, headers), nil
	}
	out, _ := json.Marshal(feedings)
	return respond(http.StatusOK, string(out), headers), nil
}

func CreateFeeding(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	var f models.Feeding
	if err := json.Unmarshal([]byte(req.Body), &f); err != nil {
		return respond(http.StatusBadRequest, `{"error":"invalid request"}`, headers), nil
	}
	f.ID = uuid.NewString()

	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	if err := s.CreateFeeding(ctx, f); err != nil {
		return respond(http.StatusInternalServerError, `{"error":"create failed"}`, headers), nil
	}
	out, _ := json.Marshal(f)
	return respond(http.StatusCreated, string(out), headers), nil
}

func UpdateFeeding(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	parts := strings.Split(req.RawPath, "/")
	id := parts[len(parts)-1]

	var f models.Feeding
	if err := json.Unmarshal([]byte(req.Body), &f); err != nil {
		return respond(http.StatusBadRequest, `{"error":"invalid request"}`, headers), nil
	}
	f.ID = id

	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	if err := s.UpdateFeeding(ctx, f); err != nil {
		return respond(http.StatusInternalServerError, `{"error":"update failed"}`, headers), nil
	}
	out, _ := json.Marshal(f)
	return respond(http.StatusOK, string(out), headers), nil
}

func DeleteFeeding(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	parts := strings.Split(req.RawPath, "/")
	id := parts[len(parts)-1]

	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	if err := s.DeleteFeeding(ctx, id); err != nil {
		return respond(http.StatusInternalServerError, `{"error":"delete failed"}`, headers), nil
	}
	return respond(http.StatusNoContent, "", headers), nil
}
