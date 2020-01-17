#--------------------------------------------------------------
# main_init.tf must be not touch! because main_init.tf is auto generate file.
#--------------------------------------------------------------

#--------------------------------------------------------------
# terraform state
#--------------------------------------------------------------
terraform {
  required_version = ">= 0.12"
  backend "azurerm" {
    storage_account_name = "##STORAGE_ACCOUNT_NAME##"
    container_name = "##CONTAINER_NAME##"
    key = "##KEY##"
  }
}

#--------------------------------------------------------------
# Provider Setting
#--------------------------------------------------------------
provider "azurerm" {
  version = "##AZURERM_VERSION##"
}
provider "azuread" {
  version = "##AZURERD_VERSION##"
}

#--------------------------------------------------------------
# Get Client Config
#--------------------------------------------------------------
data "azurerm_client_config" "current" {}
