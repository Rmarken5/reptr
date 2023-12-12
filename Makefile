.Phony:
docker-build-service:
	docker build -f ./dockerfiles/service.dockerfile . -t gcr.io/small-biz-template/markenshop/reptr:latest

.Phony:
docker-run-service:
	docker run -d --rm -p 8081:8080 -e PORT=8080 \
                                 -e MONGO_URI="mongodb://host.docker.internal:27017/?directConnection=true&serverSelectionTimeoutMS=2000" \
                                 -e DB_NAME="deck" \
                                 gcr.io/small-biz-template/markenshop/reptr:latest
docker-build-mongo:
	docker build -f ./dockerfiles/mongo.dockerfile -t reprt-mongo:latest .

docker-run-mongo:
	docker run -d --rm -p 27017:27017 -v /home/ryan/data/db:/data/db -v /home/ryan/data/log/mongodb/mongo.log:/data/log/mongodb/mongo.log reprt-mongo:latest

test:
	go test ./...

cover:
	go test -coverprofile coverage.out ./... && go tool cover -html=coverage.out

.PHONY: gen
gen:
	find . -name "*_mock.go" -type f -delete
	go generate ./...