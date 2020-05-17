package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"stackhead.io/terraform-nginx-provider/src/nginx"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
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
	}

	p.ConfigureFunc = providerConfigure(p)
	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		availableOld, availableNew := d.GetChange("directory_available")
		enabledOld, enabledNew := d.GetChange("directory_enabled")
		config := nginx.Config{
			DirectoryAvailable:          d.Get("directory_available").(string),
			DirectoryEnabled:            d.Get("directory_enabled").(string),
			EnableSymlinks:              len(d.Get("directory_enabled").(string)) > 0,
			RegenerateResources:         d.HasChange("directory_available") || d.HasChange("directory_enabled"),
			DirectoryAvailableChanged:   d.HasChange("directory_available"),
			DirectoryAvailableChangeOld: availableOld.(string),
			DirectoryAvailableChangeNew: availableNew.(string),
			DirectoryEnabledChanged:     d.HasChange("directory_enabled"),
			DirectoryEnabledChangeOld:   enabledOld.(string),
			DirectoryEnabledChangeNew:   enabledNew.(string),
		}
		return config, nil
	}
}
