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

type GetEggsResponse struct {
	Eggs []GetEggResponse `json:"eggs"`
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
	// Extract user ID from JWT claims
	claims := request.RequestContext.Authorizer.JWT.Claims
	authenticatedUser := claims["sub"]
	if authenticatedUser == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 401,
			Body:       `{"error": "Unauthorized: user ID not found in token"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Get owner from path parameter
	owner := request.PathParameters["owner"]
	if owner == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "owner parameter is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Ensure user can only access their own secrets
	if owner != authenticatedUser {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 403,
			Body:       `{"error": "Forbidden: you can only access your own secrets"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Retrieve all eggs from DynamoDB for this owner
	eggs, err := eggRepo.GetAllEggs(ctx, owner)
	if err != nil {
		println("DynamoDB Error:", err.Error())
		errorMsg := map[string]string{
			"error":   "Failed to retrieve eggs",
			"details": err.Error(),
		}
		errorBody, _ := json.Marshal(errorMsg)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       string(errorBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Decrypt each egg
	var decryptedEggs []GetEggResponse
	for _, egg := range eggs {
		// Decrypt the data key using KMS
		decryptResp, err := kmsClient.Decrypt(ctx, &kms.DecryptInput{
			CiphertextBlob: egg.EncryptedDataKey,
		})
		if err != nil {
			println("KMS Decrypt Error for SecretID", egg.SecretID, ":", err.Error())
			// Skip this egg but continue with others
			continue
		}

		// Decrypt the ciphertext using AES-256-GCM with the plaintext data key
		plaintextBytes, err := crypto.DecryptWithAESGCM(egg.Ciphertext, decryptResp.Plaintext)
		if err != nil {
			println("AES Decrypt Error for SecretID", egg.SecretID, ":", err.Error())
			// Skip this egg but continue with others
			continue
		}

		plaintext := string(plaintextBytes)

		decryptedEggs = append(decryptedEggs, GetEggResponse{
			Owner:     egg.Owner,
			SecretID:  egg.SecretID,
			Plaintext: plaintext,
			CreatedAt: egg.CreatedAt,
		})
	}

	// Return all decrypted eggs
	response := GetEggsResponse{
		Eggs: decryptedEggs,
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
