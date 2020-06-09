#--------------------------------------------------------------
# main_init.tf must be not touch! because main_init.tf is auto generate file.
# If you want to fix it, you should be fix shell/aws/files/main.template.tf.
#--------------------------------------------------------------

#--------------------------------------------------------------
# terraform state
#--------------------------------------------------------------
terraform {
  required_version = "##REQUIRED_VERSION##"
  backend "s3" {
    bucket  = "##S3_BUCKET##"
    key     = "##S3_KEY##"
    profile = "default"    # fix for environment
    region  = "##S3_REGION##" # fix for environment
  }
}

#--------------------------------------------------------------
# Provider Setting
# access key and secret key should not use.
#--------------------------------------------------------------
provider "aws" {
  profile = "default"    # fix for environment
  region  = "##REGION##" # fix for environment
  version = "##AWS_VERSION##"
}

#--------------------------------------------------------------
# my account id/region
#--------------------------------------------------------------
data "aws_caller_identity" "self" {}
data "aws_region" "self" {}
