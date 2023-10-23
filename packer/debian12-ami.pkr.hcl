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

variable "app_dbhost" {
  type    = string
  default = "${env("APP_DBHOST")}"
}

variable "app_dbuser" {
  type    = string
  default = "${env("APP_DBUSER")}"
}

variable "app_dbpassword" {
  type    = string
  default = "${env("APP_DBPASSWORD")}"
}

variable "app_dbname" {
  type    = string
  default = "${env("APP_DBNAME")}"
}

variable "app_dbport" {
  type    = string
  default = "${env("APP_DBPORT")}"
}

variable "app_serverport" {
  type    = string
  default = "${env("APP_SERVERPORT")}"
}

variable "app_default_users" {
  type    = string
  default = "${env("APP_DEFAULT_USERS_LOC")}"
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
    environment_vars=[
      "AWS_REGION=${var.aws_region}",
      "SOURCE_AMI=${var.source_ami}",
      "SSH_USERNAME=${var.ssh_username}",
      "APP_DBHOST=${var.app_dbhost}",
      "APP_DBUSER=${var.app_dbuser}",
      "APP_DBPASSWORD=${var.app_dbpassword}",
      "APP_DBNAME=${var.app_dbname}",
      "APP_DBPORT=${var.app_dbport}",
      "APP_SERVERPORT=${var.app_serverport}",
      "APP_DEFAULT_USERS_LOC=${var.app_default_users}"
    ]
    script = "./scripts/ami-setup.sh"
  }
}