package notebook

import (
	"context"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-tir/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceImages() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"images": {
				Type: schema.TypeList,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"image_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"versions": {
						Type:     schema.TypeList,
						Elem:     &schema.Schema{Type: schema.TypeString},
						Computed: true,
					},
				}},
				Computed: true,
			},
			"active_iam": {
				Type:     schema.TypeInt,
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
		},
		ReadContext: dataSourceImagesRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func dataSourceImagesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	active_iam := d.Get("active_iam").(string)
	response, err := apiClient.GetImages(active_iam)
	if err != nil {
		return diag.Errorf("Not able to find plans")
	}
	var images []interface{}
	for _, image := range response["data"].([]interface{}) {
		imageMap, _ := image.(map[string]interface{})
		versions_res := []string{}
		for _, version := range imageMap["versions"].([]interface{}) {
			versionMap, _ := version.(map[string]interface{})
			versions_res = append(versions_res, versionMap["version"].(string))
		}
		image_res := map[string]interface{}{
			"image_name": imageMap["name"],
			"versions":   versions_res,
		}
		images = append(images, image_res)
	}
	// Setting the images data to Terraform state
	d.Set("images", images)
	d.SetId("images")
	// d.Set("category","notebooks")
	log.Println("get images", d.Get("images"))
	log.Println("cate", d.Get("category"))
	log.Println("hello")

	return diags

}

// func flattenNodes(nodes *[]models.Node) []interface{} {

// 	if nodes != nil {
// 		ois := make([]interface{}, len(*nodes), len(*nodes))

// 		for i, node := range *nodes {
// 			oi := make(map[string]interface{})
// 			oi["id"] = node.ID
// 			oi["name"] = node.Name
// 			oi["is_locked"] = node.IsLocked
// 			oi["private_ip_address"] = node.PrivateIPAddress
// 			oi["public_ip_address"] = node.PublicIPAddress
// 			oi["rescue_mode_status"] = node.RescueModeStatus
// 			oi["status"] = node.Status
// 			ois[i] = oi
// 		}

// 		return ois
// 	}
// 	return make([]interface{}, 0)
// }
