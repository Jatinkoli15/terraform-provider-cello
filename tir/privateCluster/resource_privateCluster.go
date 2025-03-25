package privateCluster

import (
	"context"
	"log"
	"math"
	"strconv"
	"strings"
	"github.com/e2eterraformprovider/terraform-provider-tir/client"
	"github.com/e2eterraformprovider/terraform-provider-tir/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePrivateCluster() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the model repository. This is a required field and must be unique.",
			},
			"nodes_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The number of kubernetes nodes you want.",
			},
			"sku_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the plan name in plan listing",
			},
			"sku_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the plan type whether hourly or committed.",
			},
			"committed_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "This is optional field to specify the number of committed days you want to opt for in case of commited sku type",
			},
			"committed_instance_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Committed Instance Policy to specify what to do with chosen committed plan after committed days, whether to renew, terminate and convert to hourly",
			},
			"currency": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is currency in which you want to make payments",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Location for resource allocation ",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time at which private cluster is created",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is your project ID of platform",
			},
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is team ID ",
			},
			"active_iam": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is for Identity Access Management. ",
			},
		},
		CreateContext: resourceCreatePrivateCluster,
		UpdateContext: resourceUpdatePrivateCluster,
		ReadContext:   resourceReadPrivateCluster,
		DeleteContext: resourceDeletePrivateCluster,
	}
}

func resourceCreatePrivateCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	payload := models.PrivateCluster{
		Name:                    d.Get("name").(string),
		NodesCount:              d.Get("nodes_count").(int),
		SKUName:                 d.Get("sku_name").(string),
		SKUType:                 d.Get("sku_type").(string),
		CommittedDays:           d.Get("committed_days").(int),
		CommittedInstancePolicy: d.Get("committed_instance_policy").(string),
		Location:                d.Get("location").(string),
		Currency:                d.Get("currency").(string),
		Category:                "private_cloud",
	}

	response, err := apiClient.NewPrivateCluster(&payload, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		// log.Println(err)
		return diag.Errorf("Some error occured while creating the private Cluster. Please check the config you have provided!! %s", "fgjabsdg")
	}
	data := response["data"].(map[string]interface{})
	d.SetId(strconv.Itoa(int(math.Round(data["id"].(float64)))))
	return nil
}

func resourceReadPrivateCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	repoID := d.Id()

	response, err := apiClient.GetRepo(repoID, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			log.Println("[INFO] Repo not found, setting ID to empty")
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding item with id: %s - %v", repoID, err)
	}
	data := response["data"].(map[string]interface{})
	sku_details := data["sku_details"].(map[string]interface{})
	specs := sku_details["specs"].(map[string]interface{})
	plan := sku_details["plan"].(map[string]interface{})
	d.Set("created_at", data["created_at"])
	d.Set("name", data["name"])
	d.Set("nodes_count", data["nodes_count"])
	d.Set("sku_name", specs["name"])
	d.Set("sku_type", plan["sku_type"])
	d.Set("committed_days", plan["committed_days"])
	d.Set("currency", plan["currency"])
	return nil
}

func resourceUpdatePrivateCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return diag.Errorf("You cannot update anything please apply terraform refresh")
}

func resourceDeletePrivateCluster(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	privateClusterID := d.Id()

	_, err := apiClient.DeletePrivateCluster(privateClusterID, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			log.Println("[INFO] Private Cluster not found, setting ID to empty")
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding item with id: %s - %v", privateClusterID, err)
	}
	d.SetId("")
	return diags
}
