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

func ListMeasurements(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	measurements, err := s.ListMeasurements(ctx)
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"query failed"}`, headers), nil
	}
	out, _ := json.Marshal(measurements)
	return respond(http.StatusOK, string(out), headers), nil
}

func CreateMeasurement(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	var m models.Measurement
	if err := json.Unmarshal([]byte(req.Body), &m); err != nil {
		return respond(http.StatusBadRequest, `{"error":"invalid request"}`, headers), nil
	}
	if m.WeightLbs == nil && m.HeightIn == nil {
		return respond(http.StatusBadRequest, `{"error":"provide at least weight or height"}`, headers), nil
	}
	m.ID = uuid.NewString()

	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	if err := s.CreateMeasurement(ctx, m); err != nil {
		return respond(http.StatusInternalServerError, `{"error":"create failed"}`, headers), nil
	}
	out, _ := json.Marshal(m)
	return respond(http.StatusCreated, string(out), headers), nil
}

func UpdateMeasurement(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	parts := strings.Split(req.RawPath, "/")
	id := parts[len(parts)-1]

	var m models.Measurement
	if err := json.Unmarshal([]byte(req.Body), &m); err != nil {
		return respond(http.StatusBadRequest, `{"error":"invalid request"}`, headers), nil
	}
	if m.WeightLbs == nil && m.HeightIn == nil {
		return respond(http.StatusBadRequest, `{"error":"provide at least weight or height"}`, headers), nil
	}
	m.ID = id

	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	if err := s.UpdateMeasurement(ctx, m); err != nil {
		return respond(http.StatusInternalServerError, `{"error":"update failed"}`, headers), nil
	}
	out, _ := json.Marshal(m)
	return respond(http.StatusOK, string(out), headers), nil
}

func DeleteMeasurement(ctx context.Context, req events.APIGatewayV2HTTPRequest, headers map[string]string) (events.APIGatewayV2HTTPResponse, error) {
	parts := strings.Split(req.RawPath, "/")
	id := parts[len(parts)-1]

	s, err := store.Get()
	if err != nil {
		return respond(http.StatusInternalServerError, `{"error":"store unavailable"}`, headers), nil
	}
	if err := s.DeleteMeasurement(ctx, id); err != nil {
		return respond(http.StatusInternalServerError, `{"error":"delete failed"}`, headers), nil
	}
	return respond(http.StatusNoContent, "", headers), nil
}
