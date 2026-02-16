package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/owenHochwald/egg-carton/cmd/actions"
)

type BreakEggResponse struct {
	Message  string `json:"message"`
	Owner    string `json:"owner"`
	SecretID string `json:"secret_id"`
}

var eggRepo actions.EggRepository

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

	// Get parameters from path
	owner := request.PathParameters["owner"]
	secretID := request.PathParameters["secretId"]

	if owner == "" || secretID == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "owner and secretId parameters are required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Ensure user can only delete their own secrets
	if owner != authenticatedUser {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 403,
			Body:       `{"error": "Forbidden: you can only delete your own secrets"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Delete the egg from DynamoDB
	if err := eggRepo.BreakEgg(ctx, owner, secretID); err != nil {
		println("DynamoDB Delete Error:", err.Error())
		errorMsg := map[string]string{
			"error":   "Failed to delete egg",
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
	response := BreakEggResponse{
		Message:  "Egg deleted successfully",
		Owner:    owner,
		SecretID: secretID,
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
