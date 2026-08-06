.PHONY: build test docker-build docker-run

build:
	CGO_ENABLED=1 go build -tags no_default_driver -o pocketbase .

test:
	CGO_ENABLED=1 go test -tags no_default_driver ./...

docker-build:
	docker build -t pocketbase-sqlite-vec:local .

docker-run:
	docker compose up --build
