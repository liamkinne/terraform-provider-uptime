package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uptime "github.com/liamkinne/terraform-provider-uptime/internal/client"
)

func resourceCheckHTTP() *schema.Resource {
	return &schema.Resource{
		Description: "A HTTP check instance.",

		CreateContext: resourceCheckHTTPCreate,
		ReadContext:   resourceCheckHTTPRead,
		UpdateContext: resourceCheckHTTPUpdate,
		DeleteContext: resourceCheckHTTPDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Display name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"contact_groups": {
				Description: "Who to direct notifications to",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"locations": {
				Description: "Geographic locations for the check to run from",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Description: "Tags to identify groups of checks",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"paused": {
				Description: "Whether the check is active or inactive",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"interval": {
				Description: "The interval in seconds for the check to run",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"address": {
				Description: "The HTTP URL to run the check against",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCheckHTTPCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*uptime.ClientWithResponses)

	check := uptime.ChecksHTTP{}

	// assemble struct
	name := d.Get("name").(string)
	check.Name = &name
	check.ContactGroups = flattenSet(d.Get("contact_groups"))
	locations := flattenSet(d.Get("locations"))
	check.Locations = &locations
	tags := flattenSet(d.Get("tags"))
	check.Tags = &tags
	isPaused := d.Get("paused").(bool)
	check.IsPaused = &isPaused
	check.MspInterval = d.Get("interval").(int)
	check.MspAddress = d.Get("address").(string)

	res, err := client.PostServiceCreateHttpWithResponse(ctx, check)
	if err != nil {
		return diag.Errorf("error creating http check: %v", err)
	}

	if res.JSON200 != nil {
		var pk int = *res.JSON200.Results.Pk
		d.SetId(fmt.Sprintf("%d", pk))
	} else {
		return diag.Errorf("create http check resource response is null")
	}

	tflog.Trace(ctx, "created a http check resource")

	return nil
}

func resourceCheckHTTPRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*uptime.ClientWithResponses)

	res, err := client.GetServiceDetailWithResponse(ctx, d.Id())
	if err != nil {
		return diag.Errorf("error reading http check: %v", err)
	}

	if res.JSON200 != nil {
		d.Set("name", res.JSON200.Name)
		d.Set("contact_groups", res.JSON200.ContactGroups)
		d.Set("locations", res.JSON200.Locations)
		d.Set("tags", res.JSON200.Tags)
		d.Set("paused", res.JSON200.IsPaused)
		d.Set("interval", res.JSON200.MspInterval)
		d.Set("address", res.JSON200.MspAddress)
	} else {
		return diag.Errorf("read http check resource response is null")
	}

	tflog.Trace(ctx, "read http check resource")

	return nil
}

func resourceCheckHTTPUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*uptime.ClientWithResponses)

	check := uptime.Checks{}

	// assemble struct
	name := d.Get("name").(string)
	check.Name = &name
	check.ContactGroups = flattenSet(d.Get("contact_groups"))
	locations := flattenSet(d.Get("locations"))
	check.Locations = &locations
	tags := flattenSet(d.Get("tags"))
	check.Tags = &tags
	isPaused := d.Get("paused").(bool)
	check.IsPaused = &isPaused
	interval := d.Get("interval").(int)
	check.MspInterval = &interval
	check.MspAddress = d.Get("address").(string)

	_, err := client.PutServiceDetail(ctx, d.Id(), check)
	if err != nil {
		return diag.Errorf("error updating http check: %v", err)
	}

	return nil
}

func resourceCheckHTTPDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*uptime.ClientWithResponses)

	client.DeleteServiceDetailWithResponse(ctx, d.Id())
	return nil
}
