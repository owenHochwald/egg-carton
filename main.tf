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