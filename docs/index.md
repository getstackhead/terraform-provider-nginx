# Nginx Provider

This provider can be used to manage Nginx configurations.

## Example Usage

### For Nginx from Ubuntu/Debian repositories

When installing Nginx from the Ubuntu/Debian repositories, it will include files in the `/etc/nginx/sites-enabled` directory.
Per default, the file will be created in `/etc/nginx/sites-available`. Resources with `enabled` property will then be symlinked to `/etc/nginx/sites-enabled`.

```hcl
# No provider configuration: using defaults

# This will create file /etc/nginx/sites-available/test.conf and symlink /etc/nginx/sites-enabled/test.conf
resource "nginx_server_block" "my-server" {
  filename = "test.conf"
  enable = true
  content = <<EOF
# content of file here
EOF
}

# This will create file /etc/nginx/sites-available/test2.conf and but no symlink
resource "nginx_server_block" "my-server2" {
  filename = "test2.conf"
  enable = false
  content = <<EOF
# content of file here
EOF
}
```

### For Nginx from the official repository

When installing Nginx from the official repository, it will include files in the `/etc/nginx/conf.d` directory.
It has no mechanism for enabling/disabling configurations, so the `enabled` setting on the resource is ignored.

```hcl
# Configure provider
provider "nginx" {
  directory_available = "/etc/nginx/conf.d"  # if not set, defaults to /etc/nginx/sites-available
  enable_symlinks = false # all resources are created in the path defined at directory_available. directory_enabled is ignored.
}

# This will create file /etc/nginx/conf.d/test.conf
resource "nginx_server_block" "my-server" {
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
* `enable_symlinks` - (Optional) Create symlink for `enabled` server_block resources. If false, all resources (regardless of `enabled`) will be created at `directory_available`. Default: true

**Note:** If you want to change the `directory_available` or `directory_enabled` after resources have already been created,
destroy the resources before that. After changing the path, recreate the resources.

## Resources

### [server_block](./resources/server_block.md)