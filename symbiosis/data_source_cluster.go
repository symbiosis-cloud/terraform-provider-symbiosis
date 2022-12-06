package symbiosis

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/symbiosis-cloud/symbiosis-go"
)

func dataSourceCluster() *schema.Resource {
	return &schema.Resource{
		Description: `
    Describes a Kubernetes cluster.
    `,
		ReadContext: dataSourceClusterRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster name.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster state.",
			},
			"kube_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kubernetes version.",
			},
			"is_highly_available": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, control plane is deployed with multiple replicas for redundancy.",
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
		},
	}
}

func dataSourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterName := d.Get("name").(string)

	log.Printf("[DEBUG] Reading cluster: %s", clusterName)

	client := meta.(*symbiosis.Client)

	cluster, err := client.Cluster.Describe(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	identity, err := client.Cluster.GetIdentity(clusterName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cluster.Name)

	d.Set("name", cluster.Name)
	d.Set("state", cluster.State)
	d.Set("endpoint", cluster.APIServerEndpoint)
	d.Set("is_highly_available", cluster.IsHighlyAvailable)
	d.Set("kube_version", cluster.KubeVersion)
	d.Set("certificate", identity.CertificatePem)
	d.Set("ca_certificate", identity.ClusterCertificateAuthorityPem)
	d.Set("private_key", identity.PrivateKeyPem)

	return nil
}
