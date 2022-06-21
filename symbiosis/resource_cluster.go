package symbiosis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/symbiosis-cloud/symbiosis-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		Description: `
    Manages Kubernetes clusters.
    `,
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		DeleteContext: resourceClusterDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster name. Changing the name forces re-creation.",
			},
			"kube_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "latest",
				Description: "Kubernetes version, see symbiosis.host for valid values or \"latest\" for the most recent supported version.",
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"wait_until_initialized": {
				Type:     schema.TypeBool,
				Default:  false,
				ForceNew: true,
				Optional: true,
			},
			"configuration": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_nginx_ingress": {
							Type:     schema.TypeBool,
							Default:  false,
							ForceNew: true,
							Optional: true,
						},
					},
				},
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster API server endpoint",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating cluster: %s", d.Get("name").(string))
	var diags diag.Diagnostics
	client := meta.(*symbiosis.Client)

	configurationInput := symbiosis.ClusterConfigurationInput{
		EnableNginxIngress: d.Get("configuration.0.enable_nginx_ingress").(bool),
	}

	input := &symbiosis.ClusterInput{
		Name:          d.Get("name").(string),
		Region:        d.Get("region").(string),
		KubeVersion:   d.Get("kube_version").(string),
		Nodes:         []symbiosis.ClusterNodeInput{},
		Configuration: configurationInput,
	}

	cluster, err := client.Cluster.Create(input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cluster.Name)

	if d.Get("wait_until_initialized").(bool) {
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			c, err := client.Cluster.Describe(cluster.Name)

			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("Error describing cluster: %s", err))
			}

			if c.State != "ACTIVE" {
				return resource.RetryableError(fmt.Errorf("expected instance to be active but was in state %s", c.State))
			}

			return nil
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting cluster: %s", d.Id())
	client := meta.(*symbiosis.Client)

	err := client.Cluster.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	if d.Get("wait_until_initialized").(bool) {
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			c, err := client.Cluster.Describe(d.Id())

			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("Error describing cluster: %s", err))
			}

			if c != nil {
				return resource.RetryableError(fmt.Errorf("expected cluster to get removed but cluster is still returned from api"))
			}

			return nil
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading cluster: %s", d.Id())
	client := meta.(*symbiosis.Client)

	cluster, err := client.Cluster.Describe(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if cluster != nil {
		d.Set("name", cluster.Name)
		d.Set("state", cluster.State)
		d.Set("endpoint", cluster.APIServerEndpoint)

	} else {
		d.SetId("")
	}

	var diags diag.Diagnostics
	return diags
}
