package models

type DiaperType string

const (
	DiaperTypeWet  DiaperType = "wet"
	DiaperTypePoop DiaperType = "poop"
	DiaperTypeBoth DiaperType = "both"
)

type Diaper struct {
	ID        string     `json:"id"        dynamodbav:"id"`
	Timestamp string     `json:"timestamp" dynamodbav:"timestamp"`
	Type      DiaperType `json:"type"      dynamodbav:"type"`
	CreatedBy string     `json:"createdBy" dynamodbav:"createdBy"`
}
