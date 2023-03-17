package symbiosis

import (
	"context"
	"fmt"
	"log"
	"strings"
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
			"is_highly_available": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     false,
				Description: "When set to true it will deploy a highly available control plane with multiple replicas for redundancy.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster state [PENDING, DELETE_IN_PROGRESS, ACTIVE, FAILED]",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster API server endpoint",
			},
			"certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"ca_certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The raw kubeconfig file.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating cluster: %s", d.Get("name").(string))

	client := meta.(*symbiosis.Client)

	input := &symbiosis.ClusterInput{
		Name:              d.Get("name").(string),
		Region:            d.Get("region").(string),
		KubeVersion:       d.Get("kube_version").(string),
		Nodes:             []symbiosis.ClusterNodePoolInput{},
		IsHighlyAvailable: d.Get("is_highly_available").(bool),
	}

	cluster, err := client.Cluster.Create(input)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cluster.Name)
	d.Set("name", cluster.Name)

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

	return resourceClusterRead(ctx, d, meta)
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting cluster: %s", d.Id())
	client := meta.(*symbiosis.Client)

	err := client.Cluster.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		c, err := client.Cluster.Describe(d.Id())

		if err != nil && !strings.Contains(err.Error(), "404") {
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

	return diags
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading cluster: %s", d.Id())
	client := meta.(*symbiosis.Client)

	cluster, err := client.Cluster.Describe(d.Id())
	if err != nil && !strings.Contains(err.Error(), "404") {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Cluster resource: %v", cluster)

	identity, err := client.Cluster.GetIdentity(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if cluster != nil {
		d.Set("name", cluster.Name)
		d.Set("state", cluster.State)
		d.Set("endpoint", cluster.APIServerEndpoint)
		d.Set("is_highly_available", cluster.IsHighlyAvailable)
		d.Set("certificate", identity.CertificatePem)
		d.Set("ca_certificate", identity.ClusterCertificateAuthorityPem)
		d.Set("private_key", identity.PrivateKeyPem)
		d.Set("kubeconfig", identity.KubeConfig)

	} else {
		d.SetId("")
	}

	var diags diag.Diagnostics
	return diags
}
