package symbiosis

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
      "symbiosis_cluster":     ResourceCluster(),
      "symbiosis_node_pool":   ResourceNodePool(),
      "symbiosis_team_member": ResourceTeamMember(),
    },
    DataSourcesMap:       map[string]*schema.Resource{},
    ConfigureContextFunc: configureContext,
  }
}

type SymbiosisClient struct {
  symbiosisApi *resty.Client
}

func configureContext(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
  c := &SymbiosisClient{}
  endpoint := d.Get("endpoint").(string)
  apiKey := d.Get("api_key").(string)
  c.symbiosisApi = resty.New().SetHostURL(endpoint).SetHeader("X-Auth-ApiKey", apiKey).SetHeader("Content-Type", "application/json").SetHeader("Accept", "application/json")

  // Verify that api key is valid and has connectivity to API gateway
  resp, err := c.symbiosisApi.R().SetError(SymbiosisApiError{}).ForceContentType("application/json").Get("rest/v1/cluster")
  if err != nil {
    return c, diag.FromErr(err)
  }
  if resp.StatusCode() != 200 {
    symbiosisErr := resp.Error().(*SymbiosisApiError)
    return c, diag.FromErr(symbiosisErr)
  }

  var diags diag.Diagnostics
  return c, diags
}
