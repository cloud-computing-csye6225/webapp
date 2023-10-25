#!/bin/bash

# Set environment variables
export DEBIAN_FRONTEND=noninteractive
export CHECKPOINT_DISABLE=1

# Update and upgrade apt-get
sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get clean -y

# Extract and move files
sudo tar -xvf /tmp/webapp.tar
sudo mv ./webapp /usr/
sudo mv ./users.csv /opt/
sudo mv ./webapp.service /etc/systemd/system/

# Create application user and set permissions
sudo groupadd csye6225
sudo useradd -m -g csye6225 webapp
sudo chmod -R 750 /usr/webapp
sudo chmod -R 740 /opt/users.csv

# Install and start Postgresql
echo 'Postgres setup'
sudo apt-get install postgresql -y

# Set up systemd for webapp
sudo systemctl daemon-reload
sudo systemctl enable webapp


