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

func ListDiapers(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	diapers, err := s.ListDiapers(ctx)
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"query failed"}`, headers), nil
	}
	out, _ := json.Marshal(diapers)
	return respond(http.StatusOK, string(out), headers), nil
}

func CreateDiaper(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	var d models.Diaper
	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return respond(http.StatusBadRequest, `{"error":"invalid request"}`, headers), nil
	}
	d.ID = uuid.NewString()

	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	if err := s.CreateDiaper(ctx, d); err != nil {
		return respond(http.StatusInternalServerError, `{"error":"create failed"}`, headers), nil
	}
	out, _ := json.Marshal(d)
	return respond(http.StatusCreated, string(out), headers), nil
}

func UpdateDiaper(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	parts := strings.Split(req.RawPath, "/")
	id := parts[len(parts)-1]

	var d models.Diaper
	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return respond(http.StatusBadRequest, `{"error":"invalid request"}`, headers), nil
	}
	d.ID = id

	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	if err := s.UpdateDiaper(ctx, d); err != nil {
		return respond(http.StatusInternalServerError, `{"error":"update failed"}`, headers), nil
	}
	out, _ := json.Marshal(d)
	return respond(http.StatusOK, string(out), headers), nil
}

func DeleteDiaper(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	parts := strings.Split(req.RawPath, "/")
	id := parts[len(parts)-1]

	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	if err := s.DeleteDiaper(ctx, id); err != nil {
		return respond(http.StatusInternalServerError, `{"error":"delete failed"}`, headers), nil
	}
	return respond(http.StatusNoContent, "", headers), nil
}
