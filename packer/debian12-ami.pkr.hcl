packer {
  required_plugins {
    amazon = {
      source  = "github.com/hashicorp/amazon"
      version = ">= 1.0.0"
    }
  }
}

variable "aws_region" {
  type    = string
  default = "${env("PKR_AWS_REGION")}"
}

variable "source_ami" {
  type    = string
  default = "${env("PKR_AWS_SRC_AMI")}"
}

variable "ssh_username" {
  type    = string
  default = "${env("PKR_AWS_SSH_USERNAME")}"
}

source "amazon-ebs" "webapp-ami" {
  region          = "${var.aws_region}"
  ami_name        = "csye6225_webapp_${formatdate("YYYY_MM_DD_hh_mm_ss", timestamp())}"
  ami_description = "AMI for CSYE6225 Webapp RestAPI"
  ami_regions = [
    "${var.aws_region}",
  ]
  ami_users = [
    "089849603791",
    "080240294678",
  ]
  instance_type = "t2.micro"
  source_ami    = "${var.source_ami}"
  ssh_username  = "${var.ssh_username}"
}

build {
  sources = ["source.amazon-ebs.webapp-ami"]

  provisioner "file" {
    source      = "./builds/webapp.tar"
    destination = "/tmp/"
  }

  provisioner "shell" {
    script = "./packer/scripts/ami-setup.sh"
  }
}