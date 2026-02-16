# Simplified Cognito User Pool for Google Authentication
resource "aws_cognito_user_pool" "eggcarton_pool" {
  name = "eggcarton-user-pool"

  # Simple password policy for fallback
  password_policy {
    minimum_length                   = 12
    require_lowercase                = true
    require_uppercase                = true
    require_numbers                  = true
    require_symbols                  = true
    temporary_password_validity_days = 7
  }

  mfa_configuration = "OPTIONAL"

  software_token_mfa_configuration {
    enabled = true
  }

  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  # Email is required for Google OAuth - keep existing schema
  schema {
    name                     = "email"
    attribute_data_type      = "String"
    mutable                  = true
    required                 = true
    developer_only_attribute = false

    string_attribute_constraints {
      min_length = 5
      max_length = 255
    }
  }

  auto_verified_attributes = ["email"]

  email_configuration {
    email_sending_account = "COGNITO_DEFAULT"
  }

  user_pool_add_ons {
    advanced_security_mode = "ENFORCED"
  }

  user_attribute_update_settings {
    attributes_require_verification_before_update = ["email"]
  }

  device_configuration {
    challenge_required_on_new_device      = true
    device_only_remembered_on_user_prompt = true
  }

  tags = {
    Project = "EggCarton"
  }
}

# Cognito User Pool Domain for hosted UI
resource "aws_cognito_user_pool_domain" "eggcarton_domain" {
  domain       = "eggcarton-auth-${random_string.domain_suffix.result}"
  user_pool_id = aws_cognito_user_pool.eggcarton_pool.id
}

# Random suffix to ensure unique domain name
resource "random_string" "domain_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Google Identity Provider Configuration
resource "aws_cognito_identity_provider" "google_provider" {
  user_pool_id  = aws_cognito_user_pool.eggcarton_pool.id
  provider_name = "Google"
  provider_type = "Google"

  provider_details = {
    authorize_scopes = "email profile openid"
    client_id        = var.google_client_id
    client_secret    = var.google_client_secret
  }

  attribute_mapping = {
    email    = "email"
    username = "sub"
  }
}

# App Client for API access
resource "aws_cognito_user_pool_client" "eggcarton_client" {
  name         = "eggcarton-api-client"
  user_pool_id = aws_cognito_user_pool.eggcarton_pool.id

  depends_on = [aws_cognito_identity_provider.google_provider]

  generate_secret = false

  callback_urls = [
    "http://localhost:8080/callback",
    "https://oauth.pstmn.io/v1/callback" # For Postman/Insomnia testing
  ]

  logout_urls = [
    "http://localhost:8080/logout"
  ]

  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code", "implicit"]
  allowed_oauth_scopes = [
    "email",
    "openid",
    "profile"
  ]

  supported_identity_providers = ["Google"]

  # Token validity
  id_token_validity      = 1  # 1 hour
  access_token_validity  = 1  # 1 hour
  refresh_token_validity = 30 # 30 days

  token_validity_units {
    id_token      = "hours"
    access_token  = "hours"
    refresh_token = "days"
  }

  explicit_auth_flows = [
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH"
  ]
}

# Outputs for authentication
output "cognito_user_pool_id" {
  description = "Cognito User Pool ID"
  value       = aws_cognito_user_pool.eggcarton_pool.id
}

output "cognito_user_pool_arn" {
  description = "Cognito User Pool ARN"
  value       = aws_cognito_user_pool.eggcarton_pool.arn
}

output "cognito_app_client_id" {
  description = "Cognito App Client ID"
  value       = aws_cognito_user_pool_client.eggcarton_client.id
}

output "cognito_domain" {
  description = "Cognito hosted UI domain"
  value       = aws_cognito_user_pool_domain.eggcarton_domain.domain
}

output "google_login_url" {
  description = "Direct URL to login with Google"
  value       = "https://${aws_cognito_user_pool_domain.eggcarton_domain.domain}.auth.${var.aws_region}.amazoncognito.com/oauth2/authorize?client_id=${aws_cognito_user_pool_client.eggcarton_client.id}&response_type=token&scope=email+openid+profile&redirect_uri=https://oauth.pstmn.io/v1/callback"
}
