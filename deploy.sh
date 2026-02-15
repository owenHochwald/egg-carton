#!/bin/bash
set -e

echo "Building Lambda functions..."
./build-lambda.sh

echo ""
echo "Initializing Terraform..."
terraform init

echo ""
echo "Planning Terraform deployment..."
terraform plan -out=tfplan

echo ""
echo "Applying Terraform configuration..."
terraform apply tfplan

echo ""
echo "Deployment complete!"
echo ""
echo "Getting API endpoint..."
terraform output api_endpoint
