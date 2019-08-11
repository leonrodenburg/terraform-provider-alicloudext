# Alibaba Cloud Terraform provider extensions

Extensions to the official [Alibaba Cloud Terraform provider](https://github.com/terraform-providers/terraform-provider-alicloud/). Provides several missing resources or more complex workflows, like automatic creation of Let's Encrypt certificates & binding of API Gateway custom domains.

## Installation

To install this provider, download the latest binary for your platform. Put the executable in the correct plugin directory:

- Windows: `%APPDATA%\terraform.d\plugins`
- macOS: `~/.terraform.d/plugins`
- Linux: `~/.terraform.d/plugins`

Once the plugin is put there, you should be able to run `terraform init` after you define any of the resources of this provider.

## Configuration

There are two ways to point the `alicloudext` provider to your Alibaba Cloud credentials. Firstly, you can define them as environment variables while running Terraform commands:

- `ALICLOUD_ACCESS_KEY`: Your access key ID
- `ALICLOUD_SECRET_KEY`: Your secret key
- `ALICLOUD_REGION`: Region to deploy resources in

Alternatively, you can wire the credentials directly into the provider as variables:

```hcl
variable "access_key" {}
variable "secret_key" {}
variable "region" {}

provider alicloudext {
  access_key = var.access_key
  secret_key = var.secret_key
  region = var.region
}
```

Then define the values of the variables in a `variables.tfvars` file and make sure to not include it in version control.

## Resources

A list will soon be published here.