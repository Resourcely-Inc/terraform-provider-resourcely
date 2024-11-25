resource "resourcely_global_value" "amis" {
  name        = "EC2 AMI"
  key         = "amis"
  description = "AMI selectors."

  type = "PRESET_VALUE_OBJECT"
  options = [
    {
      key   = "ubuntu",
      label = "Ubuntu",
      value = {
        name      = "ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*",
        virt_type = "hvm",
        owner     = "099720109477",
      }
    },
    {
      key   = "centos",
      label = "CentOS",
      value = {
        name      = "CentOS Stream 9 x84_64 *",
        virt_type = "hvm",
        owner     = "125523088429",
      },
    },
    {
      key   = "debian",
      label = "Debian",
      value = {
        name      = "debian-12-amd64-*",
        virt_type = "hvm",
        owner     = "136693071363",
      },
    },
  ]
}

resource "resourcely_blueprint" "simple_instance" {
  name        = "A simple EC2 Instance"
  description = "Creates a simple EC2 instance using one of three preselected AMIs."

  cloud_provider = "PROVIDER_AMAZON"
  categories     = ["BLUEPRINT_COMPUTE"]

  is_published = true

  content = <<-EOT
    ---
    constants:
      __name: "{{ bucket }}_{{ __guid }}"
    variables:
      ami:
        description: The Linux distribution to run on this instance.
        global_value: amis
    ---
    data "aws_ami" "{{ __name }}" {
      most_recent = true

      filter {
        name   = "name"
        values = [{{ ami.name }}]
      }

      filter {
        name   = "virtualization-type"
        values = [{{ ami.virt_type }}]
      }

      owners = [{{ ami.owner }}]
    }

    resource "aws_instance" "{{ __name }}" {
      ami           = data.aws_ami.{{ __name }}.id
      instance_type = "t3.micro"

      tags = {
        Name = "HelloWorld"
      }
    }
  EOT
}
