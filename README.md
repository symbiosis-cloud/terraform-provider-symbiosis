# Terraform Provider for Symbiosis

- Symbiosis website: https://symbiosis.host
- Provider Documentation: https://registry.terraform.io/providers/symbiosis-cloud/symbiosis/latest/docs

Run the following command to build the provider

```shell
go build -o terraform-provider-symbiosis
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```
