package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCheckTag(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCheckTag,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"uptime_checktag.bar", "name", regexp.MustCompile("^ba")),
					resource.TestMatchResourceAttr(
						"uptime_checktag.bar", "color_hex", regexp.MustCompile("^#(?:[0-9a-fA-F]{3}){1,2}$")),
				),
			},
		},
	})
}

const testAccResourceCheckTag = `
resource "uptime_checktag" "bar" {
  name = "bar"
  color_hex = "#123123"
}
`
