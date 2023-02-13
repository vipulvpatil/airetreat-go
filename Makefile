protos: protos/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/*.proto

dev:
	go run .

test:
	go test ./...

coverage:
	go test -cover ./...

run: 
	docker run -i -t --rm -p 9000:9000 --env-file .env.docker socialminego

build:
	docker build -t socialminego:latest .

local-test:
	grpcurl -import-path ./protos -proto server.proto -H 'requesting_user_email: $(TEST_USER_EMAIL)' -d '{"test": "data"}' -cert certs/local/client.crt -key certs/local/client.key -cacert certs/local/SocialMineCA.crt -servername localhost 0.0.0.0:9000 protos.AiRetreatGo/Test

local-test-no-tls:
	grpcurl -import-path ./protos -proto server.proto -plaintext -H 'requesting_user_email: $(TEST_USER_EMAIL)' -d '{"test": "data"}' 0.0.0.0:9000 protos.AiRetreatGo/Test

remote-test:
	grpcurl -import-path ./protos -proto server.proto -H 'requesting_user_email: $(TEST_USER_EMAIL)' -d '{"test": "data"}' -cert certs/remote/client.crt -key certs/remote/client.key -cacert certs/remote/SocialMineCA.crt api.socialmine.io:9000 protos.AiRetreatGo/Test

remote-test-no-tls:
	grpcurl -import-path ./protos -proto server.proto -plaintext -H 'requesting_user_email: $(TEST_USER_EMAIL)' -d '{"test": "data"}' 137.184.177.206:9000 protos.AiRetreatGo/Test

psql:
	psql "$(DB_URL)"

coverage-html:
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	open coverage/coverage.html

jobs-monitoring-view:
	workwebui -redis="redis://0.0.0.0:6379" -ns="socialmine_go" -listen=":5040"
