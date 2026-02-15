package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/owenHochwald/egg-carton/cmd/actions"
	"github.com/owenHochwald/egg-carton/pkg/crypto"
)

type PutEggRequest struct {
	Owner     string `json:"owner"`
	SecretID  string `json:"secret_id"`
	Plaintext string `json:"plaintext"`
}

type PutEggResponse struct {
	Message   string `json:"message"`
	Owner     string `json:"owner"`
	SecretID  string `json:"secret_id"`
	CreatedAt string `json:"created_at"`
}

var (
	eggRepo   actions.EggRepository
	kmsClient *kms.Client
	kmsKeyID  string
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
	kmsKeyID = os.Getenv("KMS_KEY_ID")
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Log the raw request body for debugging
	println("Request Body:", request.Body)

	var req PutEggRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		// Include the actual error in the response for debugging
		errorMsg := map[string]string{
			"error":   "Invalid request body",
			"details": err.Error(),
			"body":    request.Body,
		}
		errorBody, _ := json.Marshal(errorMsg)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       string(errorBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Validate input
	if req.Owner == "" || req.SecretID == "" || req.Plaintext == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "owner, secret_id, and plaintext are required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Generate a data key using KMS
	dataKeyResp, err := kmsClient.GenerateDataKey(ctx, &kms.GenerateDataKeyInput{
		KeyId:   &kmsKeyID,
		KeySpec: "AES_256",
	})
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to generate encryption key"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Encrypt the plaintext using AES-256-GCM with the plaintext data key
	ciphertext, err := crypto.EncryptWithAESGCM([]byte(req.Plaintext), dataKeyResp.Plaintext)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       `{"error": "Failed to encrypt data"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Create the egg
	createdAt := time.Now().Format(time.RFC3339)
	egg := actions.Egg{
		Owner:            req.Owner,
		SecretID:         req.SecretID,
		Ciphertext:       ciphertext,
		EncryptedDataKey: dataKeyResp.CiphertextBlob,
		CreatedAt:        createdAt,
	}

	// Store in DynamoDB
	if err := eggRepo.PutEgg(ctx, egg); err != nil {
		// Log the actual error for debugging
		println("DynamoDB Error:", err.Error())
		errorMsg := map[string]string{
			"error":   "Failed to store egg",
			"details": err.Error(),
		}
		errorBody, _ := json.Marshal(errorMsg)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       string(errorBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Return success response
	response := PutEggResponse{
		Message:   "Egg stored successfully",
		Owner:     req.Owner,
		SecretID:  req.SecretID,
		CreatedAt: createdAt,
	}
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 201,
		Body:       string(responseBody),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(handler)
}
