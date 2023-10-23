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

# Install and start Postgresql
echo 'Postgres setup'
sudo apt-get install postgresql -y
sudo service postgresql start
sudo pg_isready

# Change PostgreSQL user password
sudo -u postgres psql -c "ALTER ROLE $APP_DBUSER WITH PASSWORD '$APP_DBPASSWORD';"
sudo service postgresql restart
sudo pg_isready

# Set environment variables for the application
echo "Setting required env variables for the application"
echo "DBHOST=$APP_DBHOST" | sudo tee -a /etc/environment
echo "DBUSER=$APP_DBUSER" | sudo tee -a /etc/environment
echo "DBPASSWORD=$APP_DBPASSWORD" | sudo tee -a /etc/environment
echo "DBNAME=$APP_DBNAME" | sudo tee -a /etc/environment
echo "DBPORT=$APP_DBPORT" | sudo tee -a /etc/environment
echo "SERVERPORT=$APP_SERVERPORT" | sudo tee -a /etc/environment
echo "DEFAULTUSERS=$APP_DEFAULT_USERS_LOC" | sudo tee -a /etc/environment

