package symbiosis

import (
	"context"

	"log"

	"github.com/symbiosis-cloud/symbiosis-go"

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
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Type of nodes for this specific pool, see docs.",
			},
			"quantity": {
				Type:        schema.TypeInt,
				Description: "Desired number of nodes for specific pool.",
				Required:    true,
			},
		},
	}
}

func resourceNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating node pool with type %v for cluster %v", d.Get("node_type").(string), d.Id())
	var diags diag.Diagnostics

	client := meta.(*symbiosis.Client)

	input := &symbiosis.NodePoolInput{
		ClusterName:  d.Get("cluster").(string),
		NodeTypeName: d.Get("node_type").(string),
		Quantity:     d.Get("quantity").(int),
	}

	resp, err := client.NodePool.Create(input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	return diags
}

func resourceNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating node pool: %s", d.Id())
	var diags diag.Diagnostics
	client := meta.(*symbiosis.Client)

	input := &symbiosis.NodePoolUpdateInput{
		Quantity: d.Get("quantity").(int),
	}
	err := client.NodePool.Update(d.Id(), input)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting node pool: %s", d.Id())
	client := meta.(*symbiosis.Client)

	err := client.NodePool.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading node pool: %s", d.Id())
	client := meta.(*symbiosis.Client)
	nodePool, err := client.NodePool.Describe(d.Id())
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
