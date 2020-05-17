package main

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
				Description: "Folder where enabled nginx configurations are stored. Set to empty string to disable symlinking",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"nginx_vhost": resourceVhost(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := nginx.Config{
		DirectoryAvailable: d.Get("directory_available").(string),
		DirectoryEnabled:   d.Get("directory_enabled").(string),
		EnableSymlinks:     len(d.Get("directory_enabled").(string)) > 0,
	}

	return config, nil
}
