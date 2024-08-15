terraform {
  required_providers {
    azureipam = {
      source = "hashicorp.com/edu/azureipam"
    }
  }
}

provider "azureipam" {
  host      = "https://ipam-xpmctiprtdfam.azurewebsites.net"
  client_id = "57696d12-0de3-46a4-8e5d-1ccf180764c0"
}


data "azureipam_admins" "example" {
}
resource "azureipam_reservation" "name" {
  space         = "test"
  block         = "test"
  smallest_cidr = true
}
output "names" {
  value = data.azureipam_admins.example

}

output "reservation" {
  value = azureipam_reservation.name

}
