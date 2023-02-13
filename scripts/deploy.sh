#!/bin/zsh

# var definitions
SSH_ADDR="root@$1"

# build docker image locally
docker build -t socialminego:latest .

# stop all existing docker containers
ssh $SSH_ADDR "docker ps -aq | xargs docker stop | xargs docker rm"

# send docker image to server and load it.
docker save socialminego:latest | bzip2 | pv | ssh $SSH_ADDR docker load

# run newly uploaded docker image
ssh $SSH_ADDR "docker run -i -t --rm -d -p 9000:9000 --env-file .env socialminego"
