package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uptime "github.com/liamkinne/terraform-provider-uptime/internal/client"
)

func resourceCheckTag() *schema.Resource {
	return &schema.Resource{
		Description: "Tag which can be applied to checks.",

		CreateContext: resourceCheckTagCreate,
		ReadContext:   resourceCheckTagRead,
		UpdateContext: resourceCheckTagUpdate,
		DeleteContext: resourceCheckTagDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Tag display name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"color_hex": {
				Description: "Tag display color as a hexadecimal value (e.g. `#012ABC`)",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCheckTagCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*uptime.ClientWithResponses)

	tag := uptime.CheckTag{}

	if name, ok := d.GetOk("name"); ok {
		tag.Tag = name.(string)
	}

	if colorHex, ok := d.GetOk("color_hex"); ok {
		tag.ColorHex = colorHex.(string)
	}

	res, err := client.PostServicetaglistWithResponse(ctx, tag)
	if err != nil {
		return diag.Errorf("error creating check tag: %v", err)
	}

	if res.JSON200 != nil {
		var pk int = *res.JSON200.Results.Pk
		d.SetId(fmt.Sprintf("%d", pk))
	} else {
		return diag.Errorf("create check tag resource response is null")
	}
	
	tflog.Trace(ctx, "created a check tag resource")

	return nil
}

func resourceCheckTagRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*uptime.ClientWithResponses)

	res, err := client.GetServiceTagDetailWithResponse(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error reading check tag: %v", err)
	}

	if res.JSON200 != nil {
		d.Set("name", res.JSON200.Tag)
		d.Set("color_hex", res.JSON200.ColorHex)
	} else {
		return diag.Errorf("read check tag resource response is null")
	}

	tflog.Trace(ctx, "read check tag resource")

	return nil
}

func resourceCheckTagUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	client := meta.(*uptime.ClientWithResponses)

	tag := uptime.CheckTag{}

	if name, ok := d.GetOk("name"); ok {
		tag.Tag = name.(string)
	}

	if colorHex, ok := d.GetOk("color_hex"); ok {
		tag.ColorHex = colorHex.(string)
	}

	_, err := client.PutServiceTagDetailWithResponse(ctx, d.Id(), tag)
	if err != nil {
		return diag.Errorf("error updating check tag: %v", err)
	}

	return nil
}

func resourceCheckTagDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*uptime.ClientWithResponses)

	client.DeleteServiceTagDetailWithResponse(ctx, d.Id())

	return nil
}
