package illumiocore

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var providerDSCCL *schema.Provider

func TestAccIllumioCCL_Read(t *testing.T) {
	listAttr := map[string]interface{}{}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactoriesInternal(&providerDSCCL),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIllumioCCLDataSourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIllumioDataSourceCCLExists("data.illumio-core_container_clusters.test", listAttr),
					testAccCheckIllumioListDataSourceSize(listAttr, "5"),
				),
			},
		},
	})
}

func testAccCheckIllumioCCLDataSourceConfig_basic() string {
	return `
	data "illumio-core_container_clusters" "test" {
		max_results = "5"
	}
	`
}

func testAccCheckIllumioDataSourceCCLExists(name string, listAttr map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("List of Container Clusters %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID was not set")
		}

		listAttr["length"] = rs.Primary.Attributes["items.#"]

		return nil
	}
}