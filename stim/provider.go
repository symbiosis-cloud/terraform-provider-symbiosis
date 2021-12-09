package stim

import (
	"context"

	"github.com/go-resty/resty/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("STIM_API_KEY", nil),
				Description: "The ApiKey used to authenticate requests towards Stim.",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("STIM_ENDPOINT", "https://api.stim.dev"),
				Description: "Endpoint for reaching the stim API. Used for debugging or when accessed through a proxy.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"stim_cluster":     ResourceCluster(),
			"stim_team_member": ResourceTeamMember(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: configureContext,
	}
}

type StimClient struct {
	stimApi *resty.Client
}

type NodePool struct {
  NodeTypeName string
  Quantity int
}

type Cluster struct {
  Name string
  State string
  NodePools []NodePool
}

type TeamMember struct {
  Email string
	Role string
}

type TeamMemberInvitation struct {
  email string
	role string
}

func configureContext(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	c := &StimClient{}
	endpoint := d.Get("endpoint").(string)
	apiKey := d.Get("api_key").(string)
	c.stimApi = resty.New().SetHostURL(endpoint).SetHeader("X-Auth-ApiKey", apiKey).SetHeader("Content-Type", "application/json").SetHeader("Accept", "application/json")

  // Verify that api key is valid and has connectivity to API gateway
  resp, err := c.stimApi.R().SetError(StimApiError{}).ForceContentType("application/json").Get("rest/v1/cluster")
  if err != nil {
    return c, diag.FromErr(err)
  }
  if resp.StatusCode() != 200 {
		stimErr := resp.Error().(*StimApiError)
    return c, diag.FromErr(stimErr)
  }

	var diags diag.Diagnostics
	return c, diags
}
