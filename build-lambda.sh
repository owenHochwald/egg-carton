#!/bin/bash
set -e

echo "Installing dependencies..."
go get github.com/aws/aws-lambda-go/lambda
go get github.com/aws/aws-lambda-go/events
go get github.com/aws/aws-sdk-go-v2/service/kms
go mod tidy

echo "Creating lambda directory..."
mkdir -p lambda

echo "Building Lambda functions..."

cd cmd/lambda/put_egg
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
zip ../../../lambda/put_egg.zip bootstrap
rm bootstrap
cd ../../..

cd cmd/lambda/get_egg
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
zip ../../../lambda/get_egg.zip bootstrap
rm bootstrap
cd ../../..

cd cmd/lambda/break_egg
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go
zip ../../../lambda/break_egg.zip bootstrap
rm bootstrap
cd ../../..

echo "Lambda functions built successfully!"