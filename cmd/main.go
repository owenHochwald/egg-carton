package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/owenHochwald/egg-carton/cmd/actions"
)

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-1"))
	if err != nil {
		log.Fatal(err)
	}

	// Create a DynamoDB client
	dynamoClient := dynamodb.NewFromConfig(cfg)

	// Initialize the EggRepository with your table name
	eggRepo := actions.NewEggRepository(dynamoClient, "EggCarton-Eggs")

	// Example: Create and store a new egg
	newEgg := actions.Egg{
		Owner:            "USER#example",
		SecretID:         "SECRET#DEMO_KEY",
		Ciphertext:       []byte("encrypted-data-here"),
		EncryptedDataKey: []byte("encrypted-key-here"),
		CreatedAt:        time.Now().Format(time.RFC3339),
	}

	err = eggRepo.PutEgg(context.TODO(), newEgg)
	if err != nil {
		log.Printf("Failed to put egg: %v\n", err)
	} else {
		log.Println("Successfully stored egg")
	}

	// Example: Retrieve an egg
	egg, err := eggRepo.GetEgg(context.TODO(), "USER#example")
	if err != nil {
		log.Printf("Failed to get egg: %v\n", err)
	} else {
		log.Printf("Retrieved egg: %v\n", egg)
	}

	// Example: Delete an egg
	// err = eggRepo.BreakEgg(context.TODO(), "USER#example", "SECRET#DEMO_KEY")
	// if err != nil {
	// 	log.Printf("Failed to break egg: %v\n", err)
	// } else {
	// 	log.Println("Successfully broke egg")
	// }
}
