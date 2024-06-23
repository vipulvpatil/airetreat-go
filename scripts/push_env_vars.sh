#!/bin/zsh

echo $1
# var definitions
SSH_ADDR="root@$1"

# first download the existing .env_airetreat from server root
scp $SSH_ADDR:~/.env_airetreat .env.downloaded

# copy .env_airetreat to server root
scp .env_airetreat $SSH_ADDR:~/
