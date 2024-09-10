package provider

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAdminsDataSource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Read testing
            {
                Config: providerConfig + `data "azureipam_admins" "test" {}`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    // Verify number of coffees returned
                    resource.TestCheckResourceAttr("data.azureipam_admins.test", "admins.#", "2"),
                ),
            },
        },
    })
}