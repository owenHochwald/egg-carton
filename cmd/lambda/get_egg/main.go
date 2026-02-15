package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/owenHochwald/egg-carton/cmd/actions"
	"github.com/owenHochwald/egg-carton/pkg/crypto"
)

type GetEggResponse struct {
	Owner     string `json:"owner"`
	SecretID  string `json:"secret_id"`
	Plaintext string `json:"plaintext"` // Decrypted secret
	CreatedAt string `json:"created_at"`
}

var (
	eggRepo   actions.EggRepository
	kmsClient *kms.Client
)

func init() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config: " + err.Error())
	}

	// Initialize DynamoDB client and repository
	dynamoClient := dynamodb.NewFromConfig(cfg)
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		tableName = "EggCarton-Eggs"
	}
	eggRepo = actions.NewEggRepository(dynamoClient, tableName)

	// Initialize KMS client
	kmsClient = kms.NewFromConfig(cfg)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	owner := request.PathParameters["owner"]
	if owner == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "owner parameter is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Retrieve the egg from DynamoDB
	egg, err := eggRepo.GetEgg(ctx, owner)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 404,
			Body:       `{"error": "Egg not found"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Decrypt the data key using KMS
	decryptResp, err := kmsClient.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: egg.EncryptedDataKey,
	})
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to decrypt data key"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Decrypt the ciphertext using AES-256-GCM with the plaintext data key
	plaintextBytes, err := crypto.DecryptWithAESGCM(egg.Ciphertext, decryptResp.Plaintext)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to decrypt secret"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	plaintext := string(plaintextBytes)

	// Return the decrypted egg
	response := GetEggResponse{
		Owner:     egg.Owner,
		SecretID:  egg.SecretID,
		Plaintext: plaintext,
		CreatedAt: egg.CreatedAt,
	}
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(handler)
}
