package symbiosis

import (
  "context"
  "fmt"
  "log"

  "github.com/google/uuid"

  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceNodePool() *schema.Resource {
  return &schema.Resource{
    Description:   `Creates node pools for Kubernetes clusters in Symbiosis.`,
    CreateContext: resourceNodePoolCreate,
    ReadContext:   resourceNodePoolRead,
    UpdateContext: resourceNodePoolUpdate,
    DeleteContext: resourceNodePoolDelete,
    Schema: map[string]*schema.Schema{
      "id": {
        Type:        schema.TypeString,
        Computed:    true,
        Description: "ID of node pool.",
      },
      "cluster": {
        Type:        schema.TypeString,
        Required:    true,
        ForceNew:    true,
        Description: "Name of cluster to create node pool in.",
      },
      "node_type": {
        Type:     schema.TypeString,
        ForceNew: true,
        Required: true,
        Description: "Type of nodes for this specific pool, see docs.",
      },
      "quantity": {
        Type:     schema.TypeInt,
        Description: "Desired number of nodes for specific pool.",
        Required: true,
      },
    },
  }
}

type PostNodePoolInput struct {
  ClusterName  string `json:"clusterName"`
  NodeTypeName string `json:"nodeTypeName"`
  Quantity     int    `json:"quantity"`
}

type PutNodePoolInput struct {
  Quantity int `json:"quantity"`
}

type PostNodePoolPayload struct {
  Id uuid.UUID `json:"nodePoolId"`
}

func resourceNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
  log.Printf("[DEBUG] Creating node pool with type %v for cluster %v", d.Get("node_type").(string), d.Id())
  var diags diag.Diagnostics

  client := meta.(*SymbiosisClient)
  api := client.symbiosisApi

  input := &PostNodePoolInput{
    ClusterName:  d.Get("cluster").(string),
    NodeTypeName: d.Get("node_type").(string),
    Quantity:     d.Get("quantity").(int),
  }

  resp, err := api.R().SetBody(input).SetResult(PostNodePoolPayload{}).SetError(SymbiosisApiError{}).ForceContentType("application/json").Post("rest/v1/node-pool")
  if err != nil {
    return diag.FromErr(err)
  }
  if resp.StatusCode() != 200 {
    symbiosisErr := resp.Error().(*SymbiosisApiError)
    if symbiosisErr.Message != "" {
      return diag.FromErr(symbiosisErr)
    }
    return append(diags, diag.Diagnostic{
      Severity: diag.Error,
      Summary:  "Unable to create node pool",
      Detail:   "Unknown error",
    })
  }
  res := resp.Result().(*PostNodePoolPayload)
  d.SetId(res.Id.String())

  return diags
}

func resourceNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
  log.Printf("[DEBUG] Updating node pool: %s", d.Id())
  var diags diag.Diagnostics
  client := meta.(*SymbiosisClient).symbiosisApi
  input := PutNodePoolInput{
    Quantity: d.Get("quantity").(int),
  }
  resp, err := client.R().SetError(SymbiosisApiError{}).SetBody(input).ForceContentType("application/json").Put(fmt.Sprintf("rest/v1/cluster/%v/node-pool/%s", d.Get("cluster").(string), d.Id()))
  if err != nil {
    return diag.FromErr(err)
  }
  if resp.StatusCode() != 200 {
    symbiosisErr := resp.Error().(*SymbiosisApiError)
    if symbiosisErr.Message != "" {
      return diag.FromErr(symbiosisErr)
    }
    return append(diags, diag.Diagnostic{
      Severity: diag.Error,
      Summary:  "Unable to update node pool",
      Detail:   "Unknown error",
    })
  }
  return diags
}

func resourceNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
  log.Printf("[DEBUG] Deleting node pool: %s", d.Id())
  client := meta.(*SymbiosisClient)
  api := client.symbiosisApi
  resp, err := api.R().SetError(SymbiosisApiError{}).ForceContentType("application/json").Delete(fmt.Sprintf("rest/v1/node-pool/%v", d.Id()))
  if err != nil {
    return diag.FromErr(err)
  }
  if resp.StatusCode() != 200 {
    symbiosisErr := resp.Error().(*SymbiosisApiError)
    if symbiosisErr.Message != "" {
      return diag.FromErr(symbiosisErr)
    }
    return diag.FromErr(err)
  }

  var diags diag.Diagnostics
  return diags
}

func resourceNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
  log.Printf("[DEBUG] Reading node pool: %s", d.Id())
  client := meta.(*SymbiosisClient)
  nodePool, err := client.describeNodePool(d.Id())
  if err != nil {
    return diag.FromErr(err)
  }
  if nodePool != nil {
    d.Set("cluster", nodePool.ClusterName)
    d.Set("node_type", nodePool.NodeTypeName)
    d.Set("quantity", nodePool.DesiredQuantity)
  } else {
    d.SetId("")
  }

  var diags diag.Diagnostics
  return diags
}
