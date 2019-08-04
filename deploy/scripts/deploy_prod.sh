#! /bin/sh

cd ~/backend

# Update our repo
git pull

# Make sure pgdata exists, otherwise docker might fail while trying to mount it.
sudo mkdir -p /var/lib/pgsql/data

# Make and set permissions on acme.json for Let's Encrypt.
# This stores our SSL certs (for HTTPS) so we can reuse the same ones even after
# restarting our docker containers.  Without this we would have to request a
# new cert every time and would eventually hit Let's Encrypt's weekly limit.
sudo mkdir -p /opt/traefik
sudo touch /opt/traefik/acme.json && sudo chmod 600 /opt/traefik/acme.json

# Set up frontend dependencies
sudo npm install -g gulp
npm install
gulp prod

# This will rebuild our serivces *then* replace the current ones with them.
# Avoiding downtime is cool
sudo docker-compose -f ~/backend/deploy/docker-compose.default.yml -f ~/backend/deploy/docker-compose.prod.yml up -d --build --force-recreate
