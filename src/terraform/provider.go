package terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"stackhead.io/terraform-nginx-provider/src/nginx"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"directory_available": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/etc/nginx/sites-available",
				Description: "Folder where all nginx configurations are stored",
			},
			"directory_enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/etc/nginx/sites-enabled",
				Description: "Folder where enabled nginx configurations are stored. Not in use if enable_symlinks=false.",
			},
			"enable_symlinks": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Create symlink for enabled server_block resources. If false, all resources (regardless of enabled) will be created at directory_available.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"nginx_server_block": resourceServerBlock(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := nginx.Config{
		DirectoryAvailable: d.Get("directory_available").(string),
		DirectoryEnabled:   d.Get("directory_enabled").(string),
		EnableSymlinks:     d.Get("enable_symlinks").(bool),
	}

	return config, nil
}
