# Terraform Nginx provider

This provider can be used to manage Nginx configurations.

## Example Usage

```
# No provider configuration: using defaults

# This will create file /etc/nginx/sites-available/test.conf and symlink /etc/nginx/sites-enabled/test.conf
resource "nginx_vhost" "my-server" {
  filename = "test.conf"
  enable = true
  content = <<EOF
# content of file here
EOF
}

# This will create file /etc/nginx/sites-available/test2.conf and but no symlink
resource "nginx_vhost" "my-server2" {
  filename = "test2.conf"
  enable = false
  content = <<EOF
# content of file here
EOF
}
```

```
# Configure provider
provider "nginx" {
  directory_available = "/etc/nginx/conf.d"  # if not set, defaults to /etc/nginx/sites-available
  directory_enabled = ""  # if not set, defaults to /etc/nginx/sites-enabled
  enable_symlinks = false # all resources are created in the path defined at directory_available. directory_enabled is ignored.
}

# This will create file /etc/nginx/conf.d/test.conf
resource "nginx_vhost" "my-server" {
  filename = "test.conf"
  content = <<EOF
# content of file here
EOF
}
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html) (e.g. `alias` and `version`), the following arguments are supported in the Nginx provider block:

* `directory_available` - (Optional) Folder where all nginx configurations are stored. Default: `/etc/nginx/sites-available`
* `directory_enabled` - (Optional) Folder where enabled nginx configurations are stored. Not in use if `enable_symlinks`=false. Default: `/etc/nginx/sites-enabled`
* `enable_symlinks` - (Optional) Create symlink for `enabled` vhost resources. If false, all resources (regardless of `enabled`) will be created at `directory_available`. Default: true

The `vhost` resource has the following arguments:

* `filename` - (Required) Name of the configuration file
* `content` - (Required) Content of the configuration file
* `enable` - (Optional) Whether to enable the resource as active configuration. If symlinks were disabled in provider, this setting is ignored. Default: true

## Development

### Build Go binary file

```yaml
go build -o dist/terraform-provider-nginx ./src
```