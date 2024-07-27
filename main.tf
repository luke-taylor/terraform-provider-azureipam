terraform {
  required_providers {
    azureipam = {
      source = "hashicorp.com/edu/azureipam"
    }
  }
}

provider "azureipam" {}


data "azureipam_admins" "example" {
}

output "names" {
  value = data.azureipam_admins.example
  
}