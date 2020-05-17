package main

import (
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"stackhead.io/terraform-nginx-provider/src/nginx"
)

func resourceServerBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerBlockCreate,
		Read:   resourceServerBlockRead,
		Update: resourceServerBlockUpdate,
		Delete: resourceServerBlockDelete,

		Schema: map[string]*schema.Schema{
			"filename": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the configuration file",
			},
			"content": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Content of the configuration file",
			},
			"enable": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable the resource as active configuration. If symlinks were disabled in provider, this setting is ignored.",
			},
		},
	}
}

func resourceServerBlockCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(nginx.Config)

	// Create file
	content := d.Get("content").(string)
	fullPathAvailable, err := nginx.CreateOrUpdateServerBlock(d.Get("filename").(string), content, config)
	if err != nil {
		return err
	}

	if config.EnableSymlinks && d.Get("enable").(bool) {
		if err := nginx.EnableServerBlock(d.Get("filename").(string), config); err != nil {
			return err
		}
	}

	d.SetId(fullPathAvailable)
	return resourceServerBlockRead(d, m)
}

func resourceServerBlockRead(d *schema.ResourceData, m interface{}) error {
	config := m.(nginx.Config)
	availablePath := config.DirectoryAvailable
	enabledPath := config.DirectoryEnabled
	content, err := nginx.ReadFile(d.Id())
	if err != nil {
		return err
	}
	d.Set("filename", filepath.Base(d.Id()))
	d.Set("content", content)
	fullEnabledPath := strings.Replace(d.Id(), availablePath, enabledPath, 1)
	d.Set("enable", nginx.FileExists(fullEnabledPath))
	return nil
}

func resourceServerBlockUpdate(d *schema.ResourceData, m interface{}) error {
	// Content changed: replace old file content
	if d.HasChange("content") {
		_, err := nginx.CreateOrUpdateServerBlock(d.Id(), d.Get("content").(string), m.(nginx.Config))
		if err != nil {
			return err
		}
	}

	// Enable changed: set or remove symlink site-enabled -> site-available
	if d.HasChange("enable") {
		if d.Get("enable").(bool) {
			if err := nginx.EnableServerBlock(d.Get("filename").(string), m.(nginx.Config)); err != nil {
				return err
			}
		} else {
			if err := nginx.DisableServerBlock(d.Get("filename").(string), m.(nginx.Config)); err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceServerBlockDelete(d *schema.ResourceData, m interface{}) error {
	if err := nginx.RemoveNginxServerBlock(d.Get("filename").(string), m.(nginx.Config)); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
