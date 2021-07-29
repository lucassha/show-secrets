## terraform

This Terraform creates the necessary AWS resources to allow for GitHub Actions to push a new release up to S3 upon a new version being labeled.

It hosts:
* S3 bucket
* S3 bucket policy
* IAM user
* IAM user policy