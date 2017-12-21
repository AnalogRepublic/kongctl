setup:
	glide install

build:
	go build

test:
	go test ./config/... ./data/... ./kong/... ./util/... .;

kong-up:
	docker-compose up

default: build