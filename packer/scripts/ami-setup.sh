#!/bin/bash

# Set environment variables
export DEBIAN_FRONTEND=noninteractive
export CHECKPOINT_DISABLE=1

# Update and upgrade apt-get
sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get clean -y

# Create application user
sudo groupadd csye6225
sudo useradd -s /bin/false -g csye6225 -d /opt/webapp -m webapp

# Create folder structure
sudo -u webapp mkdir /opt/webapp/bin
sudo -u webapp mkdir /opt/webapp/conf
sudo -u webapp mkdir /opt/webapp/data
sudo -u webapp mkdir /opt/webapp/logs

# Extract and move files
sudo tar -xvf /tmp/webapp.tar
sudo mv ./webapp /opt/webapp/bin/
sudo mv ./users.csv /opt/webapp/data/
sudo mv ./webapp.service /etc/systemd/system/

# Set permissions
sudo chmod -R 750 /opt/webapp/bin/webapp
sudo chmod -R 740 /opt/webapp/data/users.csv

# Install and start Postgresql
echo 'Postgres setup'
sudo apt-get install postgresql -y

# Set up systemd for webapp
sudo systemctl daemon-reload
sudo systemctl enable webapp


