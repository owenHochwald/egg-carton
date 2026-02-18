package actions

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type EggActions interface {
	GetEgg(ctx context.Context, owner string) (Egg, error)
	GetAllEggs(ctx context.Context, owner string) ([]Egg, error)
	PutEgg(ctx context.Context, egg Egg) error
	BreakEgg(ctx context.Context, owner, secretID string) error
}

type EggRepository struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func NewEggRepository(dynamoDbClient *dynamodb.Client, tableName string) EggRepository {
	return EggRepository{DynamoDbClient: dynamoDbClient, TableName: tableName}
}

func (r EggRepository) GetEgg(ctx context.Context, owner string) (Egg, error) {
	var egg Egg
	params, err := attributevalue.MarshalList([]interface{}{owner})
	if err != nil {
		panic(err)
	}
	response, err := r.DynamoDbClient.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("SELECT * FROM \"%v\" WHERE Owner=?",
				r.TableName)),
		Parameters: params,
	})
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", owner, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Items[0], &egg)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}
	return egg, err
}

func (r EggRepository) GetAllEggs(ctx context.Context, owner string) ([]Egg, error) {
	var eggs []Egg
	params, err := attributevalue.MarshalList([]interface{}{owner})
	if err != nil {
		panic(err)
	}
	response, err := r.DynamoDbClient.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("SELECT * FROM \"%v\" WHERE Owner=?",
				r.TableName)),
		Parameters: params,
	})
	if err != nil {
		log.Printf("Couldn't get eggs for %v. Here's why: %v\n", owner, err)
		return eggs, err
	}

	// Unmarshal all items
	err = attributevalue.UnmarshalListOfMaps(response.Items, &eggs)
	if err != nil {
		log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		return eggs, err
	}

	return eggs, nil
}

func (r EggRepository) PutEgg(ctx context.Context, egg Egg) error {
	// Use standard PutItem instead of PartiQL for proper upsert behavior
	item, err := attributevalue.MarshalMap(egg)
	if err != nil {
		log.Printf("Couldn't marshal egg to DynamoDB item. Here's why: %v\n", err)
		return err
	}
	
	_, err = r.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("Couldn't put an item. Here's why: %v\n", err)
	}
	return err
}

func (r EggRepository) BreakEgg(ctx context.Context, owner, secretID string) error {
	params, err := attributevalue.MarshalList([]interface{}{owner, secretID})
	if err != nil {
		panic(err)
	}
	_, err = r.DynamoDbClient.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("DELETE FROM \"%v\" WHERE Owner=? AND SecretID=?",
				r.TableName)),
		Parameters: params,
	})
	if err != nil {
		log.Printf("Couldn't delete that egg from the table. Here's why: %v\n", err)
	}
	return err
}
