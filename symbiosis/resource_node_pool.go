package symbiosis

import (
	"context"

	"log"

	"github.com/symbiosis-cloud/symbiosis-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceNodePool() *schema.Resource {

	resourceSchema := map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of node pool.",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of node pool",
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
		"labels": {
			Type:        schema.TypeMap,
			Description: "Node labels to be applied to the nodes",
			Optional:    true,
			ForceNew:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"taint": {
			Type:        schema.TypeSet,
			Description: "Node taints to be applied to the nodes",
			Optional:    true,
			ForceNew:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
					"effect": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Taint effect. Can be either NoSchedule, PreferNoSchedule or NoExecute. See: https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/",
						ValidateFunc: validation.StringInSlice([]string{
							string(symbiosis.EFFECT_NO_SCHEDULE),
							string(symbiosis.EFFECT_NO_EXECUTE),
							string(symbiosis.EFFECT_PREFER_NO_SCHEDULE),
						}, false),
					},
				},
			},
		},
		"autoscaling": {
			Type:     schema.TypeSet,
			ForceNew: false,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Required: true,
					},
					"min_size": {
						Type:         schema.TypeInt,
						ValidateFunc: validation.IntAtLeast(2),
						Required:     true,
					},
					"max_size": {
						Type:         schema.TypeInt,
						ValidateFunc: validation.IntAtMost(100),
						Required:     true,
					},
				},
			},
		},
	}

	return &schema.Resource{
		Description:   `Creates node pools for Kubernetes clusters in Symbiosis.`,
		CreateContext: resourceNodePoolCreate,
		ReadContext:   resourceNodePoolRead,
		UpdateContext: resourceNodePoolUpdate,
		DeleteContext: resourceNodePoolDelete,
		Schema:        resourceSchema,
	}
}

func resourceNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating node pool with type %v for cluster %v", d.Get("node_type").(string), d.Id())
	var diags diag.Diagnostics

	client := meta.(*symbiosis.Client)

	labels := expandLabels(d.Get("labels").(map[string]interface{}))
	taints := expandTaints(d.Get("taint").(*schema.Set).List())

	autoscaling := expandAutoscalingSettings(d.Get("autoscaling").(*schema.Set).List())

	input := &symbiosis.NodePoolInput{
		Name:         d.Get("name").(string),
		ClusterName:  d.Get("cluster").(string),
		NodeTypeName: d.Get("node_type").(string),
		Quantity:     d.Get("quantity").(int),
		Labels:       labels,
		Taints:       taints,
		Autoscaling:  autoscaling,
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

	autoscaling := expandAutoscalingSettings(d.Get("autoscaling").(*schema.Set).List())

	log.Printf("[DEBUG] Updating node pool: %v", autoscaling)

	input := &symbiosis.NodePoolUpdateInput{
		Quantity:    d.Get("quantity").(int),
		Autoscaling: autoscaling,
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

		d.Set("name", nodePool.Name)
		d.Set("cluster", nodePool.ClusterName)
		d.Set("node_type", nodePool.NodeTypeName)
		d.Set("quantity", nodePool.DesiredQuantity)
		d.Set("labels", flattenLabels(nodePool.Labels))
		d.Set("taints", flattenedTaints(nodePool.Taints))
		d.Set("autoscaling", flattenAutoscalingSettings(nodePool.Autoscaling))
	} else {
		d.SetId("")
	}

	var diags diag.Diagnostics
	return diags
}

func expandTaints(taints []interface{}) []symbiosis.NodeTaint {
	convertedTaints := make([]symbiosis.NodeTaint, 0, len(taints))
	for _, taint := range taints {
		input := taint.(map[string]interface{})

		t := symbiosis.NodeTaint{
			Key:    input["key"].(string),
			Value:  input["value"].(string),
			Effect: symbiosis.SchedulerEffect(input["effect"].(string)),
		}

		convertedTaints = append(convertedTaints, t)
	}

	return convertedTaints
}

func expandLabels(labels map[string]interface{}) []symbiosis.NodeLabel {
	convertedLabels := make([]symbiosis.NodeLabel, 0, len(labels))

	for key, value := range labels {
		newLabel := &symbiosis.NodeLabel{
			Key:   key,
			Value: value.(string),
		}

		convertedLabels = append(convertedLabels, *newLabel)
	}

	return convertedLabels
}

func flattenLabels(labels []*symbiosis.NodeLabel) map[string]interface{} {
	flattenedLabels := make(map[string]interface{})
	for _, label := range labels {
		flattenedLabels[label.Key] = label.Value
	}
	return flattenedLabels
}

func flattenedTaints(input []*symbiosis.NodeTaint) []interface{} {
	taints := make([]interface{}, 0)
	if input == nil {
		return taints
	}

	for _, taint := range input {
		rawTaint := map[string]interface{}{
			"key":    taint.Key,
			"value":  taint.Value,
			"effect": taint.Effect,
		}

		taints = append(taints, rawTaint)
	}

	return taints
}

func flattenAutoscalingSettings(input symbiosis.AutoscalingSettings) []interface{} {
	settings := make([]interface{}, 0)

	setting := map[string]interface{}{
		"enabled":  input.Enabled,
		"min_size": input.MinSize,
		"max_size": input.MaxSize,
	}

	settings = append(settings, setting)

	return settings
}
func expandAutoscalingSettings(settings []interface{}) symbiosis.AutoscalingSettings {
	if len(settings) == 0 {
		return symbiosis.AutoscalingSettings{Enabled: false, MinSize: 0, MaxSize: 0}
	}

	settingsMap := settings[0].(map[string]interface{})
	return symbiosis.AutoscalingSettings{
		Enabled: settingsMap["enabled"].(bool),
		MinSize: settingsMap["min_size"].(int),
		MaxSize: settingsMap["max_size"].(int),
	}
}
