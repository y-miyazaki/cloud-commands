#--------------------------------------------------------------
# main_init.tf must be not touch! because main_init.tf is auto generate file.
#--------------------------------------------------------------

#--------------------------------------------------------------
# terraform state
#--------------------------------------------------------------
terraform {
  required_version = ">= 0.12"
}

#--------------------------------------------------------------
# Provider Setting
#--------------------------------------------------------------
provider "github" {
  token        = "##GITHUB_TOKEN##"
  individual   = "##GITHUB_INDIVIDUAL##"
  ##GITHUB_ORGANIZATION##
}
