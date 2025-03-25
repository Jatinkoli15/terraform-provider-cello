package notebook

import (
	"context"
	"log"
	"reflect"
	"github.com/e2eterraformprovider/terraform-provider-tir/client"
	"github.com/e2eterraformprovider/terraform-provider-tir/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSKUPlans() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"image_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plans": {
				Type: schema.TypeList,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"name": { //
						Type:     schema.TypeString,
						Computed: true,
					},
					"cpu": { ///
						Type:     schema.TypeString,
						Computed: true,
					},
					"gpu": { //
						Type:     schema.TypeString,
						Computed: true,
					},
					"sku_type": { //
						Type:     schema.TypeString,
						Computed: true,
					},
					"unit_price": { //
						Type:     schema.TypeFloat,
						Computed: true,
					},
					"committed_days": { //
						Type:     schema.TypeInt,
						Computed: true,
					},
					"currency": { //
						Type:     schema.TypeString,
						Computed: true,
					},
					"memory": { //
						Type:     schema.TypeString,
						Computed: true,
					},
				}},
				Computed: true,
			},
			"active_iam": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "notebook",
			},
			"is_jupyterlab_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "pre-built",
			},
		},
		ReadContext: dataSourcePlansRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func dataSourcePlansRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//https://api-loki.e2enetworks.net/myaccount/api/v1/gpu/gpu_service/sku/?service=notebook&image_version_id=64&apikey=646e17f7-cb66-45e9-9bef-246c6f85f418&active_iam=99
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	active_iam := d.Get("active_iam").(string)
	node := models.ImageDetail{
		ImageName:           d.Get("image_name").(string),
		ImageVersion:        d.Get("image_version").(string),
		IsJupyterLabEnabled: d.Get("is_jupyterlab_enabled").(bool),
		ImageType:           d.Get("image_type").(string),
	}
	response, err := apiClient.GetPlans(&node, active_iam)
	if err != nil {
		return diag.Errorf("Not able to find plans")
	}
	var plans []interface{}
	// Process the CPU data
	data := response["data"].(map[string]interface{})
	if cpuItems, ok := data["CPU"].([]interface{}); ok {
		for _, cpuItem := range cpuItems {
			cpuMap := cpuItem.(map[string]interface{})
			if plansList, ok := cpuMap["plans"].([]interface{}); ok {
				for _, planItem := range plansList {
					planMap := planItem.(map[string]interface{})
					plans = append(plans, map[string]interface{}{
						"name":           cpuMap["name"],
						"cpu":            cpuMap["cpu"],
						"gpu":            cpuMap["gpu"],
						"memory":         cpuMap["memory"],
						"sku_type":       planMap["sku_type"],
						"committed_days": planMap["committed_days"],
						"unit_price":     planMap["unit_price"],
						"currency":       planMap["currency"],
					})
				}
			}
		}
	}
	// Process the GPU data
	if gpuItems, ok := data["GPU"].([]interface{}); ok {
		for _, gpuItem := range gpuItems {
			gpuMap := gpuItem.(map[string]interface{})
			if plansList, ok := gpuMap["plans"].([]interface{}); ok {
				for _, planItem := range plansList {
					planMap := planItem.(map[string]interface{})
					plans = append(plans, map[string]interface{}{
						"name":           gpuMap["name"],
						"cpu":            gpuMap["cpu"],
						"gpu":            gpuMap["gpu"],
						"memory":         gpuMap["memory"],
						"sku_type":       planMap["sku_type"],
						"committed_days": planMap["committed_days"],
						"unit_price":     planMap["unit_price"],
						"currency":       planMap["currency"],
					})
				}
			}
		}
	}
	log.Println("here i am", plans)
	log.Println("type", reflect.TypeOf(plans))
	d.SetId("plans")
	d.Set("plans", plans)
	log.Println("d.get", d.Get("plans"))
	return diags
}
