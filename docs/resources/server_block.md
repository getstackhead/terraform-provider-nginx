# Server Block Resource

This resource represents a [server block configuration file](https://www.nginx.com/resources/wiki/start/topics/examples/server_blocks/) in Nginx configuration directories.

## Example Usage

```hcl
# This will create file /etc/nginx/sites-available/test.conf and symlink /etc/nginx/sites-enabled/test.conf
resource "nginx_server_block" "my-server" {
  filename = "test.conf"
  enable = true
  variables = {
    docker_port = docker_container.web.ports.external
  }
  content = <<EOF
# content of file here
# external docker port is: {# docker_port #}
EOF
}
```

## Argument Reference

* `filename` - (Required) Name of the configuration file
* `content` - (Required) Content of the configuration file
* `enable` - (Optional) Whether to enable the resource as active configuration. If symlinks were disabled in provider, this setting is ignored. Default: true
* `markers`- (Optional) Key-Value map. Keys specified as marker (e.g. `{# key #}`, `{~ key ~}`, `{* key *}`) will be replaced by the assigned value.