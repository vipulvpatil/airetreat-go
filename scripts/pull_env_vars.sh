#!/bin/zsh

# var definitions
SSH_ADDR="root@$1"

# copy .env_airetreat from server root
scp $SSH_ADDR:~/.env_airetreat .env.downloaded
