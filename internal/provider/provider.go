package provider

import (
	"context"
	"net/http"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	uptime "github.com/liamkinne/terraform-provider-uptime/internal/client"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("UPTIME_API_KEY", nil),
				},
				"api_url": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("UPTIME_API_URL", "https://uptime.com"),
				},
				"subaccount": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("UPTIME_SUBACCOUNT", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"scaffolding_data_source": dataSourceScaffolding(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"uptime_checktag": resourceCheckTag(),
				"uptime_check_http": resourceCheckHTTP(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, data *schema.ResourceData) (any, diag.Diagnostics) {
		customProvider := func(ctx context.Context, req *http.Request) error {
			key := data.Get("api_key").(string)
            req.Header.Set("Authorization", fmt.Sprintf("Token %s", key))

            sub := data.Get("subaccount").(string)
			if sub != "" {
				req.Header.Set("X-Subaccount", sub)
			}

			userAgent := fmt.Sprintf("Terraform/%s (+https://www.terraform.io)", meta.SDKVersion)
			req.Header.Set("User-Agent", userAgent)

			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-type", "application/json")

			if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodDelete {
				req.Header.Set("Content-type", "application/json")
			}

            return nil
        }

        url := data.Get("api_url").(string)
		client, err := uptime.NewClientWithResponses(url, uptime.WithRequestEditorFn(customProvider))
		if err != nil {
			diag.Errorf("error creating new REST client: %v", err)
		}

		return client, nil
	}
}
