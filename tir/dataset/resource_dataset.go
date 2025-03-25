package dataset

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
)

func ResourceEOS() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the EOS (Elastic Object Storage) resource. This is a required field and must be unique.",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of storage for the EOS resource. Supported values are 'new_bucket' for managed storage, 'existing_bucket' for E2E S3, and 'disk' for PVC (Persistent Volume Claim).",
			},
			"encryption_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether encryption is enabled for the EOS resource. Default is false.",
			},
			"encryption_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The type of encryption used for the EOS resource. This is required if encryption is enabled.",
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The name of the bucket associated with the EOS resource. This is computed automatically.",
			},
			"bucket_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the bucket associated with the EOS resource. This is computed automatically.",
			},
			"bucket_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint URL for accessing the bucket. This is computed automatically.",
			},
			"access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The access key for the EOS resource. This is computed automatically.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret key for the EOS resource. This is computed automatically.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the EOS resource. This is computed automatically.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the EOS resource was created. This is computed automatically.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project where the EOS resource is deployed.",
			},
			"active_iam": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IAM (Identity and Access Management) role associated with the EOS resource.",
			},
			"team_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the team that owns the EOS resource.",
			},
			"disk_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "The size of the disk (in GB) allocated for the EOS resource. This is applicable only for PVC storage type.",
			},
			"pvc_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The type of PVC (Persistent Volume Claim) used for the EOS resource. This is applicable only for PVC storage type.",
			},
		},
		CreateContext: resourceCreateDataset,
		UpdateContext: resourceUpdateDataset,
		ReadContext:   resourceReadDataset,
		DeleteContext: resourceDeleteDataset,
	}
}

func resourceCreateDataset(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	var diags diag.Diagnostics

	encryptionValue := ""
	if d.Get("encryption_type").(string) == "user_managed" {
		encryptionValue = "sse-c"
	} else if d.Get("encryption_type").(string) == "e2e_managed" {
		encryptionValue = "sse-kms"
	}
	storage_type := " "
	if d.Get("storage_type").(string) == "disk" {
		storage_type = "pvc"
	} else if d.Get("storage_type").(string) == "existing_bucket" {
		storage_type = "e2e_s3"
	} else if d.Get("storage_type").(string) == "new_bucket" {
		storage_type = "managed"
	}

	bucketName := d.Get("bucket_name").(string)

	dataset := models.Dataset{
		Name:              d.Get("name").(string),
		Encryption_Enable: false,
		StorageType:       storage_type,
	}
	if d.Get("encryption_enable").(bool) {
		dataset.Encryption_Enable = true
		dataset.Encryption_Type = &encryptionValue
	}

	if storage_type == "e2e_s3" {
		dataset.BucketName = &bucketName
	} else if storage_type == "pvc" {
		disk_size, ok := d.Get("disk_size").(int)
		if !ok {
			disk_size = 10
		}
		dataset.Pvc = &models.PVCDetails{
			DiskSize: disk_size,
			PvcType:  d.Get("pvc_type").(string),
		}
	}
	response, err := client.NewDataset(&dataset, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		return diag.Errorf("Some problem occured with the creation..please check the config %s", err)
	}
	data := response["data"].(map[string]interface{})
	datasetId := data["id"].(float64)
	datasetId = math.Round(datasetId)
	d.SetId(strconv.Itoa(int(math.Round(datasetId))))
	d.Set("status", data["status"].(string))
	d.Set("created_at", data["created_at"].(string))

	return diags
}

func resourceReadDataset(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	var diags diag.Diagnostics
	datasetId := d.Id()

	response, err := client.GetDataset(datasetId, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			log.Println("[INFO] Node not found, setting ID to empty")
			d.SetId("")
			return nil
		}
		return diag.Errorf("Some problem while fetching eos details")
	}
	data := response["data"].(map[string]interface{})
	bucket, ok := data["bucket"].(map[string]interface{})
	d.Set("encryption_type", data["encryption_type"])
	d.Set("encryption_enable", data["encryption_enable"])
	d.Set("storage_type", data["storage_type"])
	if ok {
		d.Set("bucket_name", bucket["bucket_name"])
		d.Set("bucket_url", bucket["bucket_url"])
		d.Set("bucket_endpoint", bucket["endpoint"])
	}
	access_key, ok := data["access_key"].(map[string]interface{})
	if ok {
		d.Set("access_key", access_key["access_key"])
		d.Set("secret_key", access_key["secret_key"])
	}
	d.Set("status", data["status"].(string))
	d.Set("created_at", data["created_at"].(string))
	return diags
}

func resourceUpdateDataset(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
func resourceDeleteDataset(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	apiClient := m.(*client.Client)
	datasetID := d.Id()

	_, err := apiClient.DeleteDataset(datasetID, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
