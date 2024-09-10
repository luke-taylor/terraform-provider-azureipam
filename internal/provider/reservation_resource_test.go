package provider

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccReservationResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Create and Read testing
            {
                Config: providerConfig + `
resource "azureipam_reservation" "name" {
  space          = "test"
  block          = "test"
  smallest_cidr  = false
  size           = 24
  reverse_search = true
}
`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    // Verify number of items

                ),
            },
            // ImportState testing
            // {
            //     ResourceName:      "azureipam_reservation.name",
            //     ImportState:       true,
            //     ImportStateVerify: true,
            //     // The last_updated attribute does not exist in the HashiCups
            //     // API, therefore there is no value for it during import.
            //     ImportStateVerifyIgnore: []string{"last_updated"},
            // },
            // Delete testing automatically occurs in TestCase
        },
    })
}