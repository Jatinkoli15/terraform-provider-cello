package e2e

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jatinkoli15/terraform-provider-cello/client"
	"github.com/jatinkoli15/terraform-provider-cello/tir/dataset"
	"github.com/jatinkoli15/terraform-provider-cello/tir/integration"
	"github.com/jatinkoli15/terraform-provider-cello/tir/modelEndpoint"
	"github.com/jatinkoli15/terraform-provider-cello/tir/modelRepo"
	"github.com/jatinkoli15/terraform-provider-cello/tir/notebook"
	"github.com/jatinkoli15/terraform-provider-cello/tir/privateCluster"
)

// Provider function defines the schema for authentication.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Endpoint of e2e tir platform",
				Default:     "https://api.e2enetworks.com/myaccount/api/v1/gpu",
			},
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Authentication token",
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "API Key for authentication",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cello_node":        notebook.ResourceNode(),
			"cello_eos":             dataset.ResourceEOS(),
			"cello_model_repository": modelRepo.ResourceModelRepo(),
			"cello_model_endpoint":   modelEndpoint.ResourceModel(),
			"cello_integration":     integration.ResourceModelRepo(),
			"cello_private_cluster":  privateCluster.ResourcePrivateCluster(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"cello_node_images":       notebook.DataSourceImages(),
			"cello_node_plans": notebook.DataSourceSKUPlans(),
		},
		ConfigureFunc: providerConfigure, // setup the API Client
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	api_key := d.Get("api_key").(string)
	auth_token := d.Get("auth_token").(string)
	api_endpoint := d.Get("api_endpoint").(string)
	return client.NewClient(api_key, auth_token, api_endpoint), nil
}
