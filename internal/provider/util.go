package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenSet(set interface{}) []string {
	input := set.(*schema.Set)
	tags := make([]string, len(input.List()))

	if input == nil {
		return tags
	}

	for k, v := range input.List() {
		tags[k] = v.(string)
	}

	return tags
}
