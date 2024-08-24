terraform {
  required_providers {
    azureipam = {
      source = "hashicorp.com/edu/azureipam"
    }
    # azurerm = {
    #   source = "hashicorp/azurerm"
    # }
  }
}
## create virtual network in tefraform
# provider "azurerm" {
#   features {}
# }


provider "azureipam" {
  host_url         = "https://ipam-xpmctiprtdfam.azurewebsites.net"
  engine_client_id = "57696d12-0de3-46a4-8e5d-1ccf180764c0"
}


data "azureipam_admins" "example" {
}
resource "azureipam_reservation" "name" {
  space          = "test"
  block          = "test"
  smallest_cidr  = false
  size           = 24
  reverse_search = true
  # cidr = "10.0.0.0/24"
}

# resource "azurerm_resource_group" "example" {
#   name     = "example-resources"
#   location = "West Europe"
# }

# resource "azurerm_virtual_network" "example" {
#   name                = "example-network"
#   location            = azurerm_resource_group.example.location
#   resource_group_name = azurerm_resource_group.example.name
#   address_space       = [azureipam_reservation.name.cidr]

#   tags = {
#     environment = "Production"
#   }
# }

output "names" {
  value = data.azureipam_admins.example

}

output "reservation" {
  value = azureipam_reservation.name

}
