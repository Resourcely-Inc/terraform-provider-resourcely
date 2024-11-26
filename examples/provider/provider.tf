# Pin the source and version of the provider
terraform {
  required_providers {
    resourcely = {
      source  = "Resourcely-Inc/resourcely"
      version = "~> 1.0"
    }
  }
}

# No provider configuration is required. See `Authentication and
# Configuration` below to learn how to provide your Resourcley API
# credentials.
provider "resourcely" {
}
