# this terraform creates an iam user for github actions to upload 
# to S3 upon a new release
# --
# it must be run locally w/ profile "shannon"

provider "aws" {
  region  = "us-west-2"
  profile = "shannon"
}

terraform {
  backend "s3" {
    bucket  = "shannon-terraform"
    key     = "show-secrets/terraform.tfstate"
    region  = "us-west-2"
    profile = "shannon"
  }
}

resource "aws_s3_bucket" "show_secrets_bucket" {
  bucket = "lucassha-show-secrets-releases"
  acl    = "public-read"

  tags = {
    created_by_terraform = "true"
    owner                = "shannon"
  }
}

resource "aws_iam_user" "s3" {
  name = "s3_show_secrets_user"
  path = "/system/"

  tags = {
    created_by_terraform = "true"
    owner                = "shannon"
  }
}

resource "aws_iam_user_policy" "s3_user_policy" {
  name = "s3_show_secrets_user_policy"
  user = aws_iam_user.s3.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "S3AllowToReleasesBucket",
      "Action": "s3:*",
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::${aws_s3_bucket.show_secrets_bucket.bucket}"
    }
  ]
}
EOF    
}