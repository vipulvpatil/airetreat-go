#!/bin/zsh

wait_for_service_to_become_healthy() {
  health_check_passed=false
  try_count=0
  max_try_count=60
  while [[ try_count -lt max_try_count && "$health_check_passed" = false ]];
  do
    try_count=$(($try_count+1));
    echo "waiting for " $1 " to become healthy." $try_count

    http_response=$(curl -I http://$1:8180 2>/dev/null | head -n 1 | cut -d$' ' -f2);
    if [[ http_response -eq 200 ]] ; then
        health_check_passed=true
    fi

    sleep 1;
  done

  if [[ try_count -eq max_try_count ]] ; then
    echo $1 "did not become healthy"
    exit 1
  fi
}

restart_service() {
  echo "restarting: $1"

  # var definitions
  SSH_ADDR="root@$1"

  # stop all existing docker containers
  ssh $SSH_ADDR "docker ps -aq | xargs docker stop --time=60 | xargs docker rm"

  # run newly uploaded docker image
  ssh $SSH_ADDR "docker run -i -t -d -p 9100:9100 -p 8180:8180 --restart unless-stopped --env-file .env_airetreat airetreat"

  # verify service is properly started
  wait_for_service_to_become_healthy $1
}

cleanup() {
  echo "performing cleanup on $1"
  SSH_ADDR="root@$1"
  ssh $SSH_ADDR "docker rmi $(docker images -f 'dangling=true' -q)"
}

# restart to primary
restart_service $AI_RETREAT_GO_INTERNAL_IP_PRIMARY

# wait to allow LB to find healthy primary
echo 'forced sleep for 120s'
sleep 120;

# restart to secondary
restart_service $AI_RETREAT_GO_INTERNAL_IP_SECONDARY

echo "restart successful"

cleanup $AI_RETREAT_GO_INTERNAL_IP_PRIMARY
cleanup $AI_RETREAT_GO_INTERNAL_IP_SECONDARY