.Phony:
docker-build-service: gen
	docker build -f ./dockerfiles/service.dockerfile . -t gcr.io/small-biz-template/markenshop/reptr:latest

.Phony:
docker-run-service: docker-build-service
	docker run -d --rm -p 8081:8080 \
    		-e PORT=8080 \
    		-e MONGO_URI="mongodb://host.docker.internal:27017/?directConnection=true&serverSelectionTimeoutMS=2000" \
    		-e DB_NAME="deck" \
    		-e AUTH0_AUDIENCE="$(AUTH0_AUDIENCE)" \
    		-e AUTH0_CLIENT_ID="$(AUTH0_CLIENT_ID)" \
    		-e AUTH0_CLIENT_SECRET="$(AUTH0_CLIENT_SECRET)" \
    		-e AUTH0_GRANT_TYPE="$(AUTH0_GRANT_TYPE)" \
    		-e AUTH0_ENDPOINT="$(AUTH0_ENDPOINT)" \
    		-e AUTH0_CALLBACK_URL="$(AUTH0_CALLBACK_URL)" \
    		gcr.io/small-biz-template/markenshop/reptr:latest
docker-build-mongo:
	docker build -f ./dockerfiles/mongo.dockerfile -t reprt-mongo:latest .

docker-run-mongo:
	docker run -d --rm -p 27017:27017 -v /home/ryan/data/db:/data/db -v /home/ryan/data/log/mongodb/mongo.log:/data/log/mongodb/mongo.log reprt-mongo:latest

test:
	export UPDATE_SNAPS=false && go test ./...

test-update:
	export UPDATE_SNAPS=true && go test ./...

cover:
	export UPDATE_SNAPS=false && go test -coverprofile coverage.out ./... && go tool cover -html=coverage.out

.PHONY: gen
gen:
	find . -name "*_mock.go" -type f -delete
	go generate ./...
	templ generate

.PHONY: local
local:
	ENV="local" go run ./service/cmd/server