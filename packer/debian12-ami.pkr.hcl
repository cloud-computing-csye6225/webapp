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
  default = "us-east-1"
}

variable "source_ami" {
  type    = string
  default = "ami-06db4d78cb1d3bbf9"
}

variable "ssh_username" {
  type    = string
  default = "admin"
}

variable "subnet_id" {
  type    = string
  default = "subnet-06e048d94838c572f"
}

source "amazon-ebs" "webapp-ami" {
  profile         = "dev.admin"
  region          = "${var.aws_region}"
  ami_name        = "csye6225_webapp_${formatdate("YYYY_MM_DD_hh_mm_ss", timestamp())}"
  ami_description = "AMI for CSYE6225 Webapp RestAPI"
  ami_regions = [
    "${var.aws_region}",
  ]

  aws_polling {
    delay_seconds = 120
    max_attempts  = 50
  }

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
      "echo 'Setting env'",
      "export DBHOST=localhost",
      "export DBUSER=postgres",
      "export DBPASSWORD=Root@6225",
      "export DBNAME=csye6225_db",
      "export DBPORT=5432",
      "echo 'dbhost is:'",
      "echo $DBHOST",
      "echo 'Postgres setup'",
      "sudo apt-get install postgresql -y",
      "sudo service postgresql start",
      "sudo pg_isready",
      "sudo -u postgres psql -c \"ALTER ROLE $DBUSER WITH PASSWORD '$DBPASSWORD';\"",
      "PGPASSWORD=$DBPASSWORD psql -U $DBUSER -h $DBHOST -p $DBPORT",
      "sudo service postgresql restart",
      "sudo pg_isready",
    ]
  }
}