package main

import (
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"stackhead.io/terraform-nginx-provider/src/nginx"
)

func resourceVhost() *schema.Resource {
	return &schema.Resource{
		Create: resourceVhostCreate,
		Read:   resourceVhostRead,
		Update: resourceVhostUpdate,
		Delete: resourceVhostDelete,

		Schema: map[string]*schema.Schema{
			"filename": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceVhostCreate(d *schema.ResourceData, m interface{}) error {
	// Create file
	content := d.Get("content").(string)
	fullPathAvailable, err := nginx.CreateOrUpdateVhost(d.Get("filename").(string), content, m.(nginx.Config))
	if err != nil {
		return err
	}

	if d.Get("enable").(bool) {
		if err := nginx.EnableVhost(d.Get("filename").(string), m.(nginx.Config)); err != nil {
			return err
		}
	}

	d.SetId(fullPathAvailable)
	return resourceVhostRead(d, m)
}

func resourceVhostRead(d *schema.ResourceData, m interface{}) error {
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

func resourceVhostUpdate(d *schema.ResourceData, m interface{}) error {
	// Content changed: replace old file content
	if d.HasChange("content") {
		_, err := nginx.CreateOrUpdateVhost(d.Id(), d.Get("content").(string), m.(nginx.Config))
		if err != nil {
			return err
		}
	}

	// Enable changed: set or remove symlink site-enabled -> site-available
	if d.HasChange("enable") {
		if d.Get("enable").(bool) {
			if err := nginx.EnableVhost(d.Get("filename").(string), m.(nginx.Config)); err != nil {
				return err
			}
		} else {
			if err := nginx.DisableVhost(d.Get("filename").(string), m.(nginx.Config)); err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceVhostDelete(d *schema.ResourceData, m interface{}) error {
	if err := nginx.RemoveNginxVhost(d.Get("filename").(string), m.(nginx.Config)); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
