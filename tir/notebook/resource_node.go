package notebook

import (
	"context"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-tir/models"

	"github.com/e2eterraformprovider/terraform-provider-tir/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceNode() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the node. Example: 'node-020315084646'. This is a required field and must be unique.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the image used for the node. This is typically used in the case of notebooks.",
			},
			"image_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The version of the image used for the node.",
			},
			"sku_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SKU (Stock Keeping Unit) name for the node. This defines the type of resource being deployed.",
			},
			"sku_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SKU type for the node. This defines whether the node is billed hourly or on a committed basis.",
				ValidateFunc: validation.StringInSlice([]string{
					"hourly",
					"committed",
				}, false),
			},
			"committed_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The number of days the node is committed for. This is used for billing and resource allocation.",
			},
			"currency": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The currency used for billing the node. Supported values are 'INR' and 'USD'.",
				ValidateFunc: validation.StringInSlice([]string{
					"INR",
					"USD",
				}, false),
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location where the node is created. Example: 'Delhi' or 'Mumbai'.",
				ValidateFunc: validation.StringInSlice([]string{
					"Delhi",
				}, false),
			},
			"active_iam": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IAM (Identity and Access Management) role associated with the node.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project where the node is deployed.",
			},
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the team that owns the node.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tir-cluster",
				Description: "The type of cluster the node belongs to. Default is 'tir-cluster'.",
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "The size of the disk (in GB) allocated for the node. Default is 30 GB.",
			},
			"enable_ssh": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether SSH access is enabled for the node. Default is false.",
			},
			"image_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "pre-built",
				Description: "The type of image used for the node. Default is 'pre-built'.",
			},
			"is_jupyterlab_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether JupyterLab is enabled for the node. Default is true.",
			},
			"notebook_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "new",
				Description: "The type of notebook associated with the node. Default is 'new'.",
			},
			"notebook_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The URL of the notebook associated with the node.",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "notebook",
				Description: "The category of the node. Default is 'notebook'.",
			},
			"sfs_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/mnt/sfs",
				Description: "The path for shared file storage. Default is '/mnt/sfs'.",
			},
			"add_ons": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of add-ons associated with the node.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"dataset_id_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of dataset IDs associated with the node.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"public": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of public configurations for the node.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the node. This is computed automatically.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the node was created. This is computed automatically.",
			},
			"notebook_url_at_tir": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the notebook at TIR (Tensor Inference Resource). This is computed automatically.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of instance for the node. Supported values are 'free_usage' and 'paid_usage'.",
				ValidateFunc: validation.StringInSlice([]string{
					"free_usage",
					"paid_usage",
				}, false),
			},
			"committed_instance_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The policy for committed instances. This defines how committed instances are managed and billed.",
			},
			"stop_node": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether to stop the node. Default is false.",
			},
		},
		CreateContext: resourceCreateNode,
		UpdateContext: resourceUpdateNode,
		ReadContext:   resourceReadNode,
		DeleteContext: resourceDeleteNode,
	}
}

func resourceCreateNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	check_flag_for_stop_node := d.Get("stop_node").(bool)
	if check_flag_for_stop_node {
		return diag.Errorf("you can't give the stop_node as True while creating the notebook")
	}
	node := models.NodeCreate{
		Name:                    d.Get("node_name").(string),
		ClusterType:             d.Get("cluster_type").(string),
		DiskSizeInGB:            d.Get("disk_size").(int),
		EnableSSH:               d.Get("enable_ssh").(bool),
		ImageType:               d.Get("image_type").(string),
		ImageName:               d.Get("image_name").(string),
		ImageVersion:            d.Get("image_version").(string),
		InstanceType:            d.Get("instance_type").(string),
		IsJupyterLabEnabled:     d.Get("is_jupyterlab_enabled").(bool),
		NotebookType:            d.Get("notebook_type").(string),
		NotebookURL:             d.Get("notebook_url").(string),
		SfsPath:                 d.Get("sfs_path").(string),
		SKUName:                 d.Get("sku_name").(string),
		SKUType:                 d.Get("sku_type").(string),
		AddOns:                  convertStringList(d.Get("add_ons").([]interface{})),
		DatasetIDList:           convertStringList(d.Get("dataset_id_list").([]interface{})),
		Location:                d.Get("location").(string),
		Currency:                d.Get("currency").(string),
		Category:                d.Get("category").(string),
		CommittedInstancePolicy: d.Get("committed_instance_policy").(string),
		CommittedDays:           d.Get("committed_days").(int),
		PublicSSHKeys:           convertStringList(d.Get("public").([]interface{})),
	}

	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	activeIAM := d.Get("active_iam").(string)
	log.Println("BEFORE API CALL")
	response, err := client.NewNode(&node, teamID, projectID, activeIAM)
	if err != nil {
		return diag.Errorf("Some problem occured with the creation..please check the config %s", err)
	}
	log.Println("AFTER API CALL")
	data := response["data"].(map[string]interface{})
	if nodeID, ok := data["id"].(float64); ok {
		d.SetId(strconv.Itoa(int(math.Round(nodeID))))
	} else {
		return diag.Errorf("failed to extract node ID from response")
	}
	d.Set("status", data["status"].(string))
	d.Set("created_at", data["created_at"].(string))
	return nil
}

func resourceUpdateNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	nodeID := d.Id()
	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	activeIAM := d.Get("active_iam").(string)
	flag := d.Get("stop_node").(bool)

	if d.HasChange("stop_node") {
		if d.Get("sku_type") == "committed" {
			d.Set("stop_node", false)
			return diag.Errorf("You cant stop a committed node")
		}
		_, err := apiClient.UpdateStartStopNode(nodeID, projectID, teamID, activeIAM, flag)
		if err != nil {
			d.Set("stop_node", false)
			return diag.Errorf("Not able to stop/start node")
		}
	} else if d.HasChange("sku_type") || d.HasChange("sku_name") || d.HasChange("committed_days") || d.HasChange("committed_instance_policy") {
		old_sku_type, _ := d.GetChange("sku_type")
		if old_sku_type == "committed" {
			return diag.Errorf("you cannot change plan in committed plan")
		}
		if !d.HasChange("sku_type") && d.Get("status") != "stopped" {
			return diag.Errorf("You have to stop the node first to change plan")
		}
		node := models.NodeAction{
			SKUType:                 d.Get("sku_type").(string),
			SKUName:                 d.Get("sku_name").(string),
			Location:                d.Get("location").(string),
			Currency:                d.Get("currency").(string),
			Category:                d.Get("category").(string),
			CommittedInstancePolicy: d.Get("committed_instance_policy").(string),
			CommittedDays:           d.Get("committed_days").(int),
		}

		_, err := apiClient.UpdatePlanNode(&node, projectID, teamID, activeIAM, nodeID)
		if err != nil {
			return diag.Errorf("Plan changing failed")
		}
		d.Set("stop_node", false)
	} else if d.HasChange("image_name") || d.HasChange("image_version") {
		node := models.ImageDetail{
			ImageName:           d.Get("image_name").(string),
			ImageVersion:        d.Get("image_version").(string),
			IsJupyterLabEnabled: d.Get("is_jupyterlab_enabled").(bool),
			ImageType:           d.Get("image_type").(string),
		}
		_, err := apiClient.UpdateImage(&node, projectID, teamID, activeIAM, nodeID)
		if err != nil {
			return diag.Errorf("Image Update failed")
		}
	}
	return diags
}

func resourceReadNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	nodeID := d.Id()
	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	activeIAM := d.Get("active_iam").(string)

	response, err := apiClient.GetNode(nodeID, projectID, teamID, activeIAM)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			log.Println("[INFO] Node not found, setting ID to empty")
			d.SetId("")
			return nil
		}
		log.Println("[ERROR] Error fetching node:", err)
		return diag.Errorf("Error finding item with id: %s - %v", nodeID, err)
	}
	data := response["data"].(map[string]interface{})
	image_details := data["image_details"].(map[string]interface{})
	sku_details := data["sku_details"].(map[string]interface{})
	specs := sku_details["specs"].(map[string]interface{})
	plan := sku_details["plan"].(map[string]interface{})
	d.Set("created_at", data["created_at"].(string))
	d.Set("status", data["status"].(string))
	d.Set("image_name", image_details["image_name"])
	d.Set("image_version", image_details["image_version"])
	d.Set("sku_name", specs["name"])
	d.Set("sku_type", plan["sku_type"])
	d.Set("committed_days", plan["committed_days"])
	d.Set("currency", plan["currency"])
	notebook_url, ok := data["lab_url"].(string)
	if ok {
		d.Set("notebook_url_at_tir", notebook_url)
	} else {
		d.Set("notebook_url_at_tir", nil)
	}
	if d.Get("status") == "stopped" {
		d.Set("stop_node", true)
	} else {
		d.Set("stop_node", false)
	}
	return diags
	// return nil

}

func resourceDeleteNode(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	apiClient := m.(*client.Client)
	nodeID := d.Id()
	projectID := d.Get("project_id").(string)
	teamID := d.Get("team_id").(string)
	activeIAM := d.Get("active_iam").(string)

	err := apiClient.DeleteNode(nodeID, projectID, teamID, activeIAM)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
	// return nil
}

// Helper function to convert []interface{} to []string
func convertStringList(input []interface{}) []string {
	result := make([]string, len(input))
	for i, v := range input {
		result[i] = v.(string)
	}
	return result
}
