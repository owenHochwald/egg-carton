package actions

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Attribute,Type,Purpose,Example Value
// Owner (Partition Key),S,"The ""Owner"" or ""Namespace""",USER#yourname
// SecretID (Sort Key),S,The Secret Identifier,SECRET#ANTHROPIC_KEY
// Ciphertext,B,The actual encrypted API key,[Binary Data]
// EncryptedDataKey,B,The KMS-wrapped key used for this specific secret,[Binary Data]
// CreatedAt,S,ISO Timestamp,2026-02-15T08:00:00Z

type Egg struct {
	Owner            string `dynamodbav:"Owner"`
	SecretID         string `dynamodbav:"SecretID"`
	Ciphertext       []byte `dynamodbav:"Ciphertext"`
	EncryptedDataKey []byte `dynamodbav:"EncryptedDataKey"`
	CreatedAt        string `dynamodbav:"CreatedAt"`
}

// GetKey returns the composite primary key of the egg in a format that can be
// sent to DynamoDB.
func (e Egg) GetKey() map[string]types.AttributeValue {
	owner, err := attributevalue.Marshal(e.Owner)
	if err != nil {
		panic(err)
	}
	secretID, err := attributevalue.Marshal(e.SecretID)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"Owner": owner, "SecretID": secretID}
}

// String returns the owner, secret ID, and created at timestamp of the egg.
func (e Egg) String() string {
	return fmt.Sprintf("%v\n\tOwner: %v\n\tSecret ID: %v\n\tCreated At: %v\n",
		e.Ciphertext, e.EncryptedDataKey, e.CreatedAt)
}
