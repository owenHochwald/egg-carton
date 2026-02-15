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
	Owner            string `dynamodbav:"owner"`
	SecretID         string `dynamodbav:"secret_id"`
	Ciphertext       []byte `dynamodbav:"ciphertext"`
	EncryptedDataKey []byte `dynamodbav:"encrypted_data_key"`
	CreatedAt        string `dynamodbav:"created_at"`
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
	return map[string]types.AttributeValue{"owner": owner, "secret_id": secretID}
}

// String returns the owner, secret ID, and created at timestamp of the egg.
func (e Egg) String() string {
	return fmt.Sprintf("%v\n\tOwner: %v\n\tSecret ID: %v\n\tCreated At: %v\n",
		e.Ciphertext, e.EncryptedDataKey, e.CreatedAt)
}
