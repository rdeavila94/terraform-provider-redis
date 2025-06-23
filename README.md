# Terraform Provider for Redis

A Terraform provider for managing Redis key-value pairs and configurations.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Installation

### Using Terraform CLI

```bash
terraform init
```

### Manual Installation

1. Download the latest release for your platform
2. Extract the binary to your Terraform plugins directory:
   - Linux/macOS: `~/.terraform.d/plugins/registry.terraform.io/floqast/redis/`
   - Windows: `%APPDATA%\terraform.d\plugins\registry.terraform.io\floqast\redis\`

## Usage

### Provider Configuration

```hcl
terraform {
  required_providers {
    redis = {
      source  = "floqast/redis"
      version = "~> 0.0"
    }
  }
}

provider "redis" {
  redis_url = "redis://localhost:6379/0"
}
```

### Resources

#### `redis_string`

Manages a Redis string key-value pair.

```hcl
resource "redis_string" "example" {
  key   = "my_key"
  value = "my_value"
}
```

**Arguments:**

- `key` (Required) - The Redis key to manage
- `value` (Required) - The string value to store
- `overridable` (Optional) - If true, allows overriding existing Redis keys. If false, creation will fail if the key already exists. Defaults to `false`.

**Attributes:**

- `id` - The Redis key
- `key` - The Redis key
- `value` - The stored value
- `overridable` - Whether the key can override existing values

### Data Sources

Currently, this provider does not include data sources.

## Examples

### Basic Usage

```hcl
terraform {
  required_providers {
    redis = {
      source  = "floqast/redis"
      version = "~> 0.0"
    }
  }
}

provider "redis" {
  redis_url = "redis://localhost:6379/0"
}

resource "redis_string" "app_config" {
  key   = "app:config:version"
  value = "1.0.0"
}

resource "redis_string" "user_session" {
  key   = "user:session:12345"
  value = "active"
}
```

### Overriding Existing Keys

```hcl
resource "redis_string" "existing_key" {
  key         = "existing:key"
  value       = "new_value"
  overridable = true
}
```


## Development

### Building from Source

```bash
git clone https://github.com/floqast/terraform-provider-redis
cd terraform-provider-redis
go build -o terraform-provider-redis
```

### Running Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue on GitHub or contact the maintainers.

## Changelog

### v0.0.1
- Initial release
- Basic Redis string resource support
- Traefik configuration module
- Added `overridable` parameter for managing existing keys 
