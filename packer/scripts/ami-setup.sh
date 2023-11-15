#!/bin/bash

# Set environment variables
export DEBIAN_FRONTEND=noninteractive
export CHECKPOINT_DISABLE=1

# Update and upgrade apt-get
sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get clean -y

# Install cloudwatch agent
wget -P /tmp/ https://amazoncloudwatch-agent.s3.amazonaws.com/debian/amd64/latest/amazon-cloudwatch-agent.deb
sudo dpkg -i -E /tmp/amazon-cloudwatch-agent.deb

# Create application user
sudo groupadd csye6225
sudo useradd -s /bin/false -g csye6225 -d /opt/webapp -m webapp

# Create folder structure
sudo -u webapp mkdir /opt/webapp/bin
sudo -u webapp mkdir /opt/webapp/conf
sudo -u webapp mkdir /opt/webapp/data

# Move cloudwatch config file
sudo mv ./webapp.json /opt/aws/amazon-cloudwatch-agent/etc/

#Extract and move build artifacts
sudo tar -xvf /tmp/webapp.tar
sudo mv ./webapp /opt/webapp/bin/
sudo mv ./users.csv /opt/webapp/data/
sudo mv ./webapp.service /etc/systemd/system/

# Set permissions
sudo mkdir /var/log/webapp
sudo chown webapp:csye6225 /opt/webapp/bin/webapp
sudo chmod 710 /opt/webapp/bin/webapp
sudo chown webapp:csye6225 /opt/webapp/data/users.csv
sudo chmod 740 /opt/webapp/data/users.csv
sudo chown -R webapp:csye6225 /opt/webapp
sudo chown -R webapp:csye6225 /var/log/webapp

# Set up systemd for webapp
sudo systemctl daemon-reload
sudo systemctl enable webapp


