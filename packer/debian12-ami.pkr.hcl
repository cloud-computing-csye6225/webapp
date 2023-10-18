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

variable "subnet_id" {
  type    = string
  default = "${env("PKR_AWS_SUBNET")}"
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
  subnet_id     = "${var.subnet_id}"
}

build {
  sources = ["source.amazon-ebs.webapp-ami"]

  provisioner "shell" {
    environment_vars = [
      "DEBIAN_FRONTEND=noninteractive",
      "CHECKPOINT_DISABLE=1",
    ]
    inline = [
      "echo 'Update apt-get'",
      "sudo apt-get update -y",
      "echo 'upgrading apt-get'",
      "sudo apt-get upgrade -y",
      "echo 'cleaning apt-get'",
      "sudo apt-get clean -y",
      "echo 'Postgres setup'",
      "sudo apt-get install postgresql -y",
      "sudo service postgresql start",
      "sudo pg_isready",
      "echo ${var.app_dbuser}",
      "echo ${var.app_dbpassword}",
      "sudo -u postgres psql -c \"ALTER ROLE ${var.app_dbuser} WITH PASSWORD '${var.app_dbpassword}';\"",
      "sudo service postgresql restart",
      "sudo pg_isready",
      "echo \"Setting required env variables for the application\"",
      "echo DBHOST=${var.app_dbhost} | sudo tee -a /etc/environment",
      "echo DBUSER=${var.app_dbuser} | sudo tee -a /etc/environment",
      "echo DBPASSWORD=${var.app_dbpassword} | sudo tee -a /etc/environment",
      "echo DBNAME=${var.app_dbname} | sudo tee -a /etc/environment",
      "echo DBPORT=${var.app_dbport} | sudo tee -a /etc/environment",
      "echo SERVERPORT=${var.app_serverport} | sudo tee -a /etc/environment",
      "echo DEFAULTUSERS=${var.app_default_users}| sudo tee -a /etc/environment",
    ]
  }

  provisioner "file" {
    source      = "./builds/webapp.tar"
    destination = "/tmp/"
  }

  provisioner "shell" {
    environment_vars = [
      "DEBIAN_FRONTEND=noninteractive",
      "CHECKPOINT_DISABLE=1",
    ]
    inline = [
      "sudo tar -xvf /tmp/webapp.tar",
      "sudo mv ./webapp /usr/",
      "sudo mv ./users.csv /opt/",
    ]
  }
}