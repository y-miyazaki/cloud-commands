#--------------------------------------------------------------
# main_init.tf must be not touch! because main_init.tf is auto generate file.
#--------------------------------------------------------------

#--------------------------------------------------------------
# terraform state
#--------------------------------------------------------------
terraform {
  required_version = "##REQUIRED_VERSION##"
  backend "gcs" {
    credentials = "##GOOGLE_CLOUD_KEYFILE_JSON##"
    bucket      = "##GCS_BUCKET##"
    prefix      = "##GCS_PREFIX##"
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
  version = "##GOOGLE_VERSION##"
}
provider "google-beta" {
  credentials = "${file("##GOOGLE_CLOUD_KEYFILE_JSON##")}"
  project     = "##PROJECT_ID##"
  region      = "##REGION##"
  zone        = "##ZONE##"
  version = "##GOOGLE_BETA_VERSION##"
}
#--------------------------------------------------------------
# Information
#--------------------------------------------------------------
data "google_project" "project" {}
