package modelRepo

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

func ResourceModelRepo() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the model repository. This is a required field and must be unique.",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of storage for the model repository. Supported values are 'new_bucket' for managed storage, 'existing_bucket' for E2E S3, and 'disk' for PVC (Persistent Volume Claim).",
			},
			"model_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of model stored in the repository. This defines the category or framework of the model (e.g., TensorFlow, PyTorch).",
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The name of the bucket associated with the model repository. This is optional and will be auto-generated if not provided.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(old) != ""
				},
			},
			"bucket_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the bucket associated with the model repository. This is computed automatically.",
			},
			"bucket_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint URL for accessing the bucket. This is computed automatically.",
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The access key for the model repository. This is optional and will be auto-generated if not provided.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(old) != ""
				},
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The secret key for the model repository. This is optional and will be auto-generated if not provided.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(old) != ""
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the model repository. This is computed automatically.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the model repository was created. This is computed automatically.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project where the model repository is deployed.",
			},
			"active_iam": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IAM (Identity and Access Management) role associated with the model repository.",
			},
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the team that owns the model repository.",
			},
		},
		CreateContext: resourceCreateModelRepo,
		UpdateContext: resourceUpdateModelRepo,
		ReadContext:   resourceReadModelRepo,
		DeleteContext: resourceDeleteModelRepo,
	}
}

func resourceCreateModelRepo(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)

	storage_type := ""
	if d.Get("storage_type") == "new" {
		storage_type = "managed"
	} else if d.Get("storage_type") == "existing" {
		storage_type = "e2e_s3"
	} else if d.Get("storage_type") == "external" {
		storage_type = "external"
	}

	repo := models.ModelRepo{
		Name:        d.Get("name").(string),
		BucketName:  d.Get("bucket_name").(string),
		ModelType:   d.Get("model_type").(string),
		StorageType: storage_type,
		SecretKey:   d.Get("secret_key").(string),
		AccessKey:   d.Get("access_key").(string),
	}

	response, err := apiClient.NewRepo(&repo, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		log.Println(err)
		return diag.Errorf("Some error occured while creating the model repository. Please check the config you have provided!! %s", err)
	}
	data := response["data"].(map[string]interface{})
	d.SetId(strconv.Itoa(int(math.Round(data["id"].(float64)))))
	return nil
}

func resourceReadModelRepo(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	bucket := data["bucket"].(map[string]interface{})
	access_key := data["access_key"].(map[string]interface{})
	d.Set("bucket_name", bucket["bucket_name"])
	d.Set("bucket_url", bucket["bucket_url"])
	d.Set("bucket_endpoint", bucket["endpoint"])
	d.Set("access_key", access_key["access_key"])
	d.Set("secret_key", access_key["secret_key"])
	d.Set("status", data["status"])
	d.Set("created_at", data["created_at"])
	d.Set("model_type", data["model_type"])
	d.Set("name", data["name"])

	return nil
}

func resourceUpdateModelRepo(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceDeleteModelRepo(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	repoID := d.Id()

	_, err := apiClient.DeleteRepo(repoID, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			log.Println("[INFO] Repo not found, setting ID to empty")
			d.SetId("")
			return nil
		}
		d.SetId("")
		return diag.Errorf("Error finding item with id: %s - %v", repoID, err)
	}
	d.SetId("")
	return diags
}
