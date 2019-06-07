package stats

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tableName = "members"

type member struct {
	ID    string `dynamodbav:"id"`
	Score int    `dynamodbav:"score"`
}

type DynamoDBStatsService struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDBStatsService(db *dynamodb.DynamoDB) DynamoDBStatsService {
	return DynamoDBStatsService{
		db: db,
	}
}

func (s DynamoDBStatsService) Get() (map[string]int, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	output, err := s.db.Scan(input)
	if err != nil {
		return nil, err
	}

	members := []member{}
	if err := dynamodbattribute.UnmarshalListOfMaps(output.Items, &members); err != nil {
		return nil, err
	}

	stats := map[string]int{}
	for _, m := range members {
		stats[m.ID] = m.Score
	}

	return stats, nil
}

func (s DynamoDBStatsService) Increment(member string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(member)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#S": aws.String("score"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incr": {N: aws.String("1")},
		},
		UpdateExpression: aws.String("ADD #S :incr"),
	}

	_, err := s.db.UpdateItem(input)
	return err
}
