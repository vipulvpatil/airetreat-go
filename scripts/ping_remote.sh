if [ $# -eq 0 ];
then
  target=10
else
  target=$1
fi

echo $target

sum=0
while [[ sum -lt $target ]];
do 
grpcurl -import-path protos -proto server.proto -H 'requesting_user_email: '$TEST_USER_EMAIL -d '{"test": "data"}' -cert certs/remote/client.crt -key certs/remote/client.key -cacert certs/remote/AiRetreatCA.crt -servername $AI_RETREAT_GO_RESERVED_IP $AI_RETREAT_GO_RESERVED_IP:9000 protos.AiRetreatGo/Test;

sum=$(($sum+1));
echo $sum
sleep 0.1;

done
