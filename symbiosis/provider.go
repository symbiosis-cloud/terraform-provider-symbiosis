package symbiosis

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symbiosis-cloud/symbiosis-go"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SYMBIOSIS_API_KEY", nil),
				Description: "The ApiKey used to authenticate requests towards Symbiosis.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYMBIOSIS_ENDPOINT", "https://api.symbiosis.host"),
				Description: "Endpoint for reaching the symbiosis API. Used for debugging or when accessed through a proxy.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"symbiosis_cluster":                 ResourceCluster(),
			"symbiosis_node_pool":               ResourceNodePool(),
			"symbiosis_team_member":             ResourceTeamMember(),
			"symbiosis_cluster_service_account": ResourceClusterServiceAccount(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: configureContext,
	}
}

func configureContext(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	endpoint := d.Get("endpoint").(string)
	apiKey := d.Get("api_key").(string)

	c, err := symbiosis.NewClientFromAPIKey(apiKey, symbiosis.WithEndpoint(endpoint))

	if err != nil {
		return nil, diag.FromErr(err)
	}

	// Verify that api key is valid and has connectivity to API gateway
	clusters, err := c.Cluster.List(10, 0)
	if err != nil {
		return c, diag.FromErr(err)
	}
	if clusters == nil {
		return c, diag.FromErr(errors.New("Failed to read API result"))
	}

	var diags diag.Diagnostics
	return c, diags
}
