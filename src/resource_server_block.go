package main

import (
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"stackhead.io/terraform-nginx-provider/src/nginx"
)

func serverConfiguration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"listen": {
				Type:     schema.TypeInt,
				Default:  80,
				Optional: true,
			},
			"server_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"forward_https": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"forward_acme": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"locations": {
				Type: schema.TypeList,
				Elem: schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"auth_basic": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auth_basic_user_file": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"root": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"use_php_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Required:    true,
				ForceNew:    true,
				Description: "Server block locations",
			},
		},
	}
}

func resourceServerBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerBlockCreate,
		Read:   resourceServerBlockRead,
		Update: resourceServerBlockUpdate,
		Delete: resourceServerBlockDelete,

		Schema: map[string]*schema.Schema{
			"filename": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the configuration file",
			},
			"configurations": {
				Type:        schema.TypeList,
				Elem:        serverConfiguration(),
				Required:    true,
				ForceNew:    true,
				Description: "Server block configurations",
			},
			"markers": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Markers in content that should be replaced",
			},
			"markers_split": {
				Type:        schema.TypeMap,
				Default:     "",
				Description: "Define marker name as key and the character where the string is split as value. Chunks can be accessed as array",
				Optional:    true,
			},
			"enable": {
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
	content := nginx.BuildConfiguration(d)
	fullPathAvailable, err := nginx.CreateOrUpdateServerBlock(d.Get("filename").(string), content, config, d.Get("markers").(map[string]interface{}), d.Get("markers_split").(map[string]interface{}))
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
	// todo: parse file and set properties - use https://github.com/lytics/confl
	d.Set("content", content)
	fullEnabledPath := strings.Replace(d.Id(), availablePath, enabledPath, 1)
	d.Set("enable", nginx.FileExists(fullEnabledPath))
	return nil
}

func resourceServerBlockUpdate(d *schema.ResourceData, m interface{}) error {
	// Content changed: replace old file content
	if d.HasChange("configurations") || d.HasChange("variables") {
		_, err := nginx.CreateOrUpdateServerBlock(d.Id(), nginx.BuildConfiguration(d), m.(nginx.Config), d.Get("markers").(map[string]interface{}), d.Get("markers_split").(map[string]interface{}))
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
