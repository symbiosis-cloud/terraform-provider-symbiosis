package symbiosis

import (
	"context"
	"log"
	"time"

	"github.com/symbiosis-cloud/symbiosis-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceClusterServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description: `
    Manages Kubernetes clusters service accounts.
    `,
		CreateContext: resourceClusterServiceAccountCreate,
		ReadContext:   resourceClusterServiceAccountRead,
		DeleteContext: resourceClusterServiceAccountDelete,
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster name. Changing the name forces re-creation.",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Service account token",
			},
			"cluster_ca_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Cluster CA certificate",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceClusterServiceAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating service account: %s", d.Get("cluster_name").(string))
	var diags diag.Diagnostics
	clusterName := d.Get("cluster_name").(string)
	client := meta.(*symbiosis.Client)

	serviceaccount, err := client.Cluster.CreateServiceAccountForSelf(clusterName)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serviceaccount.ID)

	return diags
}

func resourceClusterServiceAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting service account: %s", d.Id())
	client := meta.(*symbiosis.Client)
	clusterName := d.Get("cluster_name").(string)

	err := client.Cluster.DeleteServiceAccount(clusterName, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	return diags
}

func resourceClusterServiceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading service account: %s", d.Id())
	client := meta.(*symbiosis.Client)
	clusterName := d.Get("cluster_name").(string)

	serviceAccount, err := client.Cluster.GetServiceAccount(clusterName, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if serviceAccount != nil {
		d.Set("cluster_certificate_authority", serviceAccount.ClusterCertificateAuthority)
		d.Set("kube_config", serviceAccount.KubeConfig)
		d.Set("token", serviceAccount.ServiceAccountToken)
		d.Set("cluster_ca_certificate", serviceAccount.ClusterCertificateAuthority)
	} else {
		d.SetId("")
	}

	var diags diag.Diagnostics
	return diags
}
