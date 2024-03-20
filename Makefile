.PHONY: test genmock up-build down down-volumes

test:
	go test ./...

genmock:
	mockery --keeptree --all

up:
	docker-compose up --build

down:
	docker-compose down

down-volumes:
	docker-compose down --volumes