package integration
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
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hugging_face_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active_iam": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		CreateContext: resourceCreateIntegration,
		UpdateContext: resourceUpdateIntegration,
		ReadContext:   resourceReadIntegration,
		DeleteContext: resourceDeleteIntegration,
	}
}

func resourceCreateIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	payload := models.Integration{
		IntegrationDetails: map[string]interface{}{
			"hugging_face_token": d.Get("hugging_face_token").(string),
		},
		IntegrationType: d.Get("integration_type").(string),
		Name:            d.Get("name").(string),
	}

	response, err := apiClient.NewIntegration(&payload, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		log.Println(err)
		return diag.Errorf("Some error occured while creating the model repository. Please check the config you have provided!!")
	}
	data := response["data"].(map[string]interface{})
	d.SetId(strconv.Itoa(int(math.Round(data["id"].(float64)))))
	return nil
}

func resourceReadIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceUpdateIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceDeleteIntegration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	integrationID := d.Id()

	_, err := apiClient.DeleteIntegration(integrationID, d.Get("project_id").(string), d.Get("team_id").(string), d.Get("active_iam").(string))
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			log.Println("[INFO] Repo not found, setting ID to empty")
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding item with id: %s - %v", integrationID, err)
	}
	d.SetId("")
	return diags
}
