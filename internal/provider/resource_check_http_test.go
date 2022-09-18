package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCheckHTTP(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCheckHTTP,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"uptime_check_http.foo", "name", regexp.MustCompile("^ba")),
				),
			},
		},
	})
}

const testAccResourceCheckHTTP= `
resource "uptime_check_http" "foo" {
  name = "bar"
  address = "https://example.com/"
  contact_groups = ["Test"]
  locations = ["Australia", "Singapore", "US East", "US West", "US Central", "United Kingdom", "Japan"]
  interval = 10
}
`
