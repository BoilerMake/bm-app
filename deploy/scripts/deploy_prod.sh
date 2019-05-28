#! /bin/sh

cd ~/backend

# Update our repo
git pull

# Make and set permissions on acme.json for let's encrypt
sudo mkdir -p /opt/traefik
sudo touch /opt/traefik/acme.json && sudo chmod 600 /opt/traefik/acme.json

# This will rebuild our serivces *then* replace the current ones with them.
# Avoiding downtime is cool
sudo docker-compose -f deploy/docker-compose.base.yml -f deploy/docker-compose.prod.yml up -d --build
