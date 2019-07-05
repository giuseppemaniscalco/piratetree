compose_up:
	docker-compose up -d

test:
	go test ./... -race