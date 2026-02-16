variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-west-1"
}

variable "google_client_id" {
  description = "Google OAuth 2.0 Client ID for Cognito federation"
  type        = string
  sensitive   = true
}

variable "google_client_secret" {
  description = "Google OAuth 2.0 Client Secret for Cognito federation"
  type        = string
  sensitive   = true
}

# Optional: Chrome Extension ID
variable "chrome_extension_id" {
  description = "Chrome Extension ID for callback URL configuration"
  type        = string
  default     = ""
}
