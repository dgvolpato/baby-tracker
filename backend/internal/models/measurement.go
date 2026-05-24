package models

type Measurement struct {
	ID        string   `json:"id"               dynamodbav:"id"`
	Timestamp string   `json:"timestamp"        dynamodbav:"timestamp"`
	WeightLbs *float64 `json:"weightLbs,omitempty" dynamodbav:"weightLbs,omitempty"`
	HeightIn  *float64 `json:"heightIn,omitempty"  dynamodbav:"heightIn,omitempty"`
	CreatedBy string   `json:"createdBy"        dynamodbav:"createdBy"`
}
