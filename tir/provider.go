package e2e

import (
	"github.com/e2eterraformprovider/terraform-provider-tir/client"
	"github.com/e2eterraformprovider/terraform-provider-tir/tir/dataset"
	"github.com/e2eterraformprovider/terraform-provider-tir/tir/integration"
	"github.com/e2eterraformprovider/terraform-provider-tir/tir/modelEndpoint"
	"github.com/e2eterraformprovider/terraform-provider-tir/tir/modelRepo"
	"github.com/e2eterraformprovider/terraform-provider-tir/tir/notebook"
	"github.com/e2eterraformprovider/terraform-provider-tir/tir/privateCluster"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"tir_notebook":        notebook.ResourceNode(),
			"tir_eos":             dataset.ResourceEOS(),
			"tir_modelRepository": modelRepo.ResourceModelRepo(),
			"tir_modelEndpoint":   modelEndpoint.ResourceModel(),
			"tir_integration":     integration.ResourceModelRepo(),
			"tir_privateCluster":  privateCluster.ResourcePrivateCluster(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"tir_notebook":       notebook.DataSourceImages(),
			"tir_notebook_plans": notebook.DataSourceSKUPlans(),
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
