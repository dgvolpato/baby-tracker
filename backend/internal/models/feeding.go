package models

type FeedingType string

const (
	FeedingTypeFormula FeedingType = "formula"
	FeedingTypeBreast  FeedingType = "breast"
)

type Feeding struct {
	ID        string      `json:"id"        dynamodbav:"id"`
	Timestamp string      `json:"timestamp" dynamodbav:"timestamp"`
	Type      FeedingType `json:"type"      dynamodbav:"type"`
	Oz        float64     `json:"oz"        dynamodbav:"oz"`
	CreatedBy string      `json:"createdBy" dynamodbav:"createdBy"`
}
