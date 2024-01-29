BINARY_NAME=app

build:
	go build -o bin/$(BINARY_NAME) cmd/$(BINARY_NAME)/main.go

run: build
	./bin/$(BINARY_NAME)

clean:
	rm -rf bin/$(BINARY_NAME)

test:
	go test -v ./... --cover

rundb:
	docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:16

startdb:
	docker start postgres

stopdb:
	docker stop postgres