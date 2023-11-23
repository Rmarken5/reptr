docker-build-mongo:
	docker build -f ./dockerfiles/mongo.dockerfile -t reprt-mongo:latest .

docker-run-mongo:
	docker run -d --rm -p 27017:27017 -v /home/ryan/data/db:/data/db -v /home/ryan/data/log/mongodb/mongo.log:/data/log/mongodb/mongo.log reprt-mongo:latest

test:
	go test ./...

cover:
	go test -coverprofile coverage.out ./... && go tool cover -html=coverage.out