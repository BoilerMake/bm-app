#! /bin/sh

cd ~/backend

# Update our repo
git pull

# This will rebuild our serivces *then* replace the current ones with them.
# Avoiding downtime is cool
docker-compose -f deploy/docker-compose.base.yml -f deploy/docker-compose.testing.yml up -d --build
