package main

import (
	"log"
	"github.com/e2eterraformprovider/terraform-provider-tir/tir"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return e2e.Provider()
		},
	})
}
