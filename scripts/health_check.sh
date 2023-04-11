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

# TODO: Remove next line. This next commented line should never work as the port is never exposed beyond the Load balancer
# grpcurl -- plaintext -import-path protos -proto server.proto api.airetreat.io:9090 protos.AiRetreatGoHealth/Check;

grpcurl --plaintext -import-path protos -proto server.proto localhost:9090 protos.AiRetreatGoHealth/Check;

sum=$(($sum+1));
echo $sum
sleep 0.1;

done
