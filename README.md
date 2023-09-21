# Terraform Resourcely Provider

The [Resourcly
Provider](https://registry.terraform.io/providers/Resourcely-Inc/resourcely/latest)
allows [Terraform](https://terraform.io/) to manage
[Resourcely](https://resourcely.io) blueprints, global context, etc

# Authoring Your Own Blueprints
Blueprints are templated Terraform files that streamline creating similar configurations. They use tags as placeholders, enhancing reusability and customizability across diverse setups. During resource creation, the blueprint is rendered into final Terraform configuration with tags substituted by the values provided by a developer. Resourcely will open a PR containing the ready-to-apply, customized Terraform config.

This is where you find more information on how to author your own [Blueprints](https://docs.resourcely.com/getting-started/using-resourcely/setting-up-blueprints/authoring-your-own-blueprints)
