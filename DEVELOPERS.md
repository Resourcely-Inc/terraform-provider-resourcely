## Development

### Install the Provider

The easiest way to install and build the provider is to run the following command:
```
make install
```

NOTE: The below instructions are only relevant once we have published this provider:

To instruct Terrform to use the locally built version of the provider, add the
following to your `~/.terraformrc` file.

```hcl
provider_installation {

  dev_overrides {
      "registry.terraform.io/Resourcely-Inc/resourcely" = "<GOBIN>"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Replace `<GOBIN>` with the path to the directory containing the installed
provider. Normally this is `<GOPATH>/bin`, where `<GOPATH>` is the output of
`go env GOPATH`.

### Build the Provider

```shell
go install
```

### Generate the Documentation

```shell
go generate
```

### Run the Acceptance Tests

Note: These tests create (and destroy) real resources. `RESOURCELY_AUTH_TOKEN`
must be set to a valid auth token for the tests to succeed.

```shell
make testacc
```
