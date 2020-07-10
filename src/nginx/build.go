package nginx

import (
	"bytes"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ServerConfiguration struct {
	Port         int
	ServerName   string
	UseHttps     bool
	ForwardHttps bool
	ForwardAcme  string
	Location     ServerLocation
}

type ServerLocation struct {
	Path              string
	AuthBasic         string
	AuthBasicUserFile string
	Root              string
	PhpVersion        string
}

func BuildConfiguration(d *schema.ResourceData) string {
	var tpl bytes.Buffer

	for _, c := range d.Get("configurations").([]interface{}) {
		config := c.(*schema.ResourceData)
		configuration := mapConfiguration(config)

		t, err := template.New("config").ParseFiles("templates/nginx_server_block.conf")
		if err != nil {
			panic(err)
		}
		err = t.Execute(&tpl, configuration)
		if err != nil {
			panic(err)
		}
	}

	return tpl.String()
}

func mapConfiguration(config *schema.ResourceData) ServerConfiguration {
	locationConfig := config.Get("location").(*schema.ResourceData)
	return ServerConfiguration{
		Port:         config.Get("listen").(int),
		ServerName:   config.Get("server_name").(string),
		UseHttps:     config.Get("https").(bool),
		ForwardHttps: config.Get("forward_https").(bool),
		ForwardAcme:  config.Get("forward_acme").(string),
		Location: ServerLocation{
			Path:              locationConfig.Get("path").(string),
			AuthBasic:         locationConfig.Get("auth_basic").(string),
			AuthBasicUserFile: locationConfig.Get("auth_basic_user_file").(string),
			Root:              locationConfig.Get("root").(string),
			PhpVersion:        locationConfig.Get("use_php_version").(string),
		},
	}
}
