#--------------------------------------------------------------
# main_init.tf must be not touch! because main_init.tf is auto generate file.
#--------------------------------------------------------------

#--------------------------------------------------------------
# terraform state
#--------------------------------------------------------------
terraform {
  required_version = ">= 0.12"
  backend "gcs" {
    credentials = "##GOOGLE_CLOUD_KEYFILE_JSON##"
    bucket      = "##BUCKET##"
    prefix      = "##PREFIX##"
  }
}

#--------------------------------------------------------------
# Provider Setting
#--------------------------------------------------------------
provider "google" {
  credentials = "${file("##GOOGLE_CLOUD_KEYFILE_JSON##")}"
  project     = "##PROJECT_ID##"
  region      = "##REGION##"
  zone        = "##ZONE##"
}
provider "google-beta" {
  credentials = "${file("##GOOGLE_CLOUD_KEYFILE_JSON##")}"
  project     = "##PROJECT_ID##"
  region      = "##REGION##"
  zone        = "##ZONE##"
}
#--------------------------------------------------------------
# Information
#--------------------------------------------------------------
data "google_project" "project" {}
