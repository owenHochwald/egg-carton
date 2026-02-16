provider "aws" {
  region = "us-west-1"
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_dynamodb_table" "egg_carton" {
  name         = "EggCarton-Eggs"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Owner"
  range_key    = "SecretID"

  attribute {
    name = "Owner"
    type = "S"
  }

  attribute {
    name = "SecretID"
    type = "S"
  }

  server_side_encryption {
    enabled     = true
    kms_key_arn = aws_kms_key.vault_master.arn
  }

  tags = {
    Project = "EggCarton"
  }
}

resource "aws_kms_key" "vault_master" {
  description             = "Master key for encrypting secrets"
  deletion_window_in_days = 7
  enable_key_rotation     = true

}
resource "aws_kms_alias" "vault_master_alias" {
  name          = "alias/eggcarton-master"
  target_key_id = aws_kms_key.vault_master.key_id
}

resource "aws_s3_bucket" "carton" {
  bucket = "carton-storage-bucket"

  tags = {
    Name        = "Carton Storage Bucket"
    Environment = "Dev"
  }
}

# IAM Role for Lambda Functions
resource "aws_iam_role" "lambda_exec" {
  name = "eggcarton_lambda_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

# IAM Policy for DynamoDB Access
resource "aws_iam_policy" "lambda_dynamodb_policy" {
  name        = "eggcarton_lambda_dynamodb_policy"
  description = "IAM policy for Lambda to access DynamoDB"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan",
          "dynamodb:ExecuteStatement",
          "dynamodb:PartiQLInsert",
          "dynamodb:PartiQLSelect",
          "dynamodb:PartiQLUpdate",
          "dynamodb:PartiQLDelete"
        ]
        Resource = aws_dynamodb_table.egg_carton.arn
      }
    ]
  })
}

# IAM Policy for KMS Access
resource "aws_iam_policy" "lambda_kms_policy" {
  name        = "eggcarton_lambda_kms_policy"
  description = "IAM policy for Lambda to access KMS"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "kms:Decrypt",
          "kms:Encrypt",
          "kms:GenerateDataKey",
          "kms:DescribeKey"
        ]
        Resource = aws_kms_key.vault_master.arn
      }
    ]
  })
}

# Attach policies to Lambda role
resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_dynamodb" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_dynamodb_policy.arn
}

resource "aws_iam_role_policy_attachment" "lambda_kms" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_kms_policy.arn
}

resource "aws_lambda_function" "put_egg" {
  filename      = "lambda/put_egg.zip"
  function_name = "eggcarton_put_egg"
  role          = aws_iam_role.lambda_exec.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  timeout       = 30

  source_code_hash = fileexists("lambda/put_egg.zip") ? filebase64sha256("lambda/put_egg.zip") : null

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.egg_carton.name
      KMS_KEY_ID = aws_kms_key.vault_master.key_id
    }
  }

  tags = {
    Project = "EggCarton"
  }
}

resource "aws_lambda_function" "get_egg" {
  filename      = "lambda/get_egg.zip"
  function_name = "eggcarton_get_egg"
  role          = aws_iam_role.lambda_exec.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  timeout       = 30

  source_code_hash = fileexists("lambda/get_egg.zip") ? filebase64sha256("lambda/get_egg.zip") : null

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.egg_carton.name
      KMS_KEY_ID = aws_kms_key.vault_master.key_id
    }
  }

  tags = {
    Project = "EggCarton"
  }
}

resource "aws_lambda_function" "break_egg" {
  filename      = "lambda/break_egg.zip"
  function_name = "eggcarton_break_egg"
  role          = aws_iam_role.lambda_exec.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  timeout       = 30

  source_code_hash = fileexists("lambda/break_egg.zip") ? filebase64sha256("lambda/break_egg.zip") : null

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.egg_carton.name
      KMS_KEY_ID = aws_kms_key.vault_master.key_id
    }
  }

  tags = {
    Project = "EggCarton"
  }
}

# API Gateway
resource "aws_apigatewayv2_api" "eggcarton_api" {
  name          = "eggcarton-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "POST", "PUT", "DELETE"]
    allow_headers = ["*"]
  }
}

# JWT Authorizer using Cognito
resource "aws_apigatewayv2_authorizer" "cognito" {
  api_id           = aws_apigatewayv2_api.eggcarton_api.id
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = "cognito-authorizer"

  jwt_configuration {
    audience = [aws_cognito_user_pool_client.eggcarton_client.id]
    issuer   = "https://cognito-idp.${var.aws_region}.amazonaws.com/${aws_cognito_user_pool.eggcarton_pool.id}"
  }
}

resource "aws_apigatewayv2_stage" "eggcarton_stage" {
  api_id      = aws_apigatewayv2_api.eggcarton_api.id
  name        = "dev"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gw.arn

    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
    })
  }
}

resource "aws_cloudwatch_log_group" "api_gw" {
  name              = "/aws/api_gw/eggcarton"
  retention_in_days = 30
}

# API Gateway Integrations
resource "aws_apigatewayv2_integration" "put_egg" {
  api_id                 = aws_apigatewayv2_api.eggcarton_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.put_egg.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_integration" "get_egg" {
  api_id                 = aws_apigatewayv2_api.eggcarton_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.get_egg.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_integration" "break_egg" {
  api_id                 = aws_apigatewayv2_api.eggcarton_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.break_egg.invoke_arn
  payload_format_version = "2.0"
}

# API Gateway Routes with Cognito Authorization
resource "aws_apigatewayv2_route" "put_egg" {
  api_id             = aws_apigatewayv2_api.eggcarton_api.id
  route_key          = "POST /eggs"
  target             = "integrations/${aws_apigatewayv2_integration.put_egg.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "get_egg" {
  api_id             = aws_apigatewayv2_api.eggcarton_api.id
  route_key          = "GET /eggs/{owner}"
  target             = "integrations/${aws_apigatewayv2_integration.get_egg.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

resource "aws_apigatewayv2_route" "break_egg" {
  api_id             = aws_apigatewayv2_api.eggcarton_api.id
  route_key          = "DELETE /eggs/{owner}/{secretId}"
  target             = "integrations/${aws_apigatewayv2_integration.break_egg.id}"
  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
}

# Lambda Permissions for API Gateway
resource "aws_lambda_permission" "put_egg" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.put_egg.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.eggcarton_api.execution_arn}/*/*"
}

resource "aws_lambda_permission" "get_egg" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_egg.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.eggcarton_api.execution_arn}/*/*"
}

resource "aws_lambda_permission" "break_egg" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.break_egg.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.eggcarton_api.execution_arn}/*/*"
}

# Outputs
output "api_endpoint" {
  description = "API Gateway endpoint URL"
  value       = aws_apigatewayv2_stage.eggcarton_stage.invoke_url
}

output "kms_key_id" {
  description = "KMS Key ID for encryption"
  value       = aws_kms_key.vault_master.key_id
}

output "dynamodb_table_name" {
  description = "DynamoDB table name"
  value       = aws_dynamodb_table.egg_carton.name
}