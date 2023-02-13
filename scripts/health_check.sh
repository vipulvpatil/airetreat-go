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
grpcurl -import-path protos -proto server.proto  -cert certs/remote/client.crt -key certs/remote/client.key -cacert certs/remote/AiRetreatCA.crt api.airetreat.io:9090 protos.AiRetreatGoHealth/Check;

# grpcurl -import-path protos -proto server.proto  -cert certs/local/client.crt -key certs/local/client.key -cacert certs/local/AiRetreatCA.crt localhost:9090 protos.AiRetreatGoHealth/Check;

sum=$(($sum+1));
echo $sum
sleep 0.1;

done
