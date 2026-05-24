package store

import (
	"context"
	"os"
	"sort"

	"baby-tracker/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const babyPK = "BABY#default"

type Store struct {
	client *dynamodb.Client
	table  string
}

var instance *Store

func Get() (*Store, error) {
	if instance != nil {
		return instance, nil
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	instance = &Store{
		client: dynamodb.NewFromConfig(cfg),
		table:  os.Getenv("TABLE_NAME"),
	}
	return instance, nil
}

type dbItem struct {
	PK        string  `dynamodbav:"pk"`
	SK        string  `dynamodbav:"sk"`
	ID        string  `dynamodbav:"id"`
	Timestamp string  `dynamodbav:"timestamp"`
	Type      string  `dynamodbav:"type"`
	Oz        float64 `dynamodbav:"oz"`
	CreatedBy string  `dynamodbav:"createdBy"`
}

func (s *Store) ListFeedings(ctx context.Context) ([]models.Feeding, error) {
	out, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.table),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: babyPK},
			":prefix": &types.AttributeValueMemberS{Value: "FEEDING#"},
		},
	})
	if err != nil {
		return nil, err
	}

	var items []dbItem
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &items); err != nil {
		return nil, err
	}

	feedings := make([]models.Feeding, len(items))
	for i, it := range items {
		feedings[i] = models.Feeding{
			ID:        it.ID,
			Timestamp: it.Timestamp,
			Type:      models.FeedingType(it.Type),
			Oz:        it.Oz,
			CreatedBy: it.CreatedBy,
		}
	}
	sort.Slice(feedings, func(i, j int) bool {
		return feedings[i].Timestamp > feedings[j].Timestamp
	})
	return feedings, nil
}

func (s *Store) CreateFeeding(ctx context.Context, f models.Feeding) error {
	it := dbItem{
		PK:        babyPK,
		SK:        "FEEDING#" + f.ID,
		ID:        f.ID,
		Timestamp: f.Timestamp,
		Type:      string(f.Type),
		Oz:        f.Oz,
		CreatedBy: f.CreatedBy,
	}
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return err
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.table),
		Item:      av,
	})
	return err
}

func (s *Store) UpdateFeeding(ctx context.Context, f models.Feeding) error {
	tsVal, _ := attributevalue.Marshal(f.Timestamp)
	typeVal, _ := attributevalue.Marshal(string(f.Type))
	ozVal, _ := attributevalue.Marshal(f.Oz)

	_, err := s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(s.table),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: babyPK},
			"sk": &types.AttributeValueMemberS{Value: "FEEDING#" + f.ID},
		},
		UpdateExpression: aws.String("SET #ts = :ts, #type = :type, oz = :oz"),
		ExpressionAttributeNames: map[string]string{
			"#ts":   "timestamp",
			"#type": "type",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ts":   tsVal,
			":type": typeVal,
			":oz":   ozVal,
		},
	})
	return err
}

func (s *Store) DeleteFeeding(ctx context.Context, id string) error {
	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.table),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: babyPK},
			"sk": &types.AttributeValueMemberS{Value: "FEEDING#" + id},
		},
	})
	return err
}
