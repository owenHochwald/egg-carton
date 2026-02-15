terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.92"
    }
  }

  backend "s3" {
    bucket         = "carton-bucket-state-0002"
    key            = "dev/egg-vault.tfstate"
    region         = "us-west-1"
    dynamodb_table = "terraform-lock"
    encrypt        = true
  }


  required_version = ">= 1.2"
}
