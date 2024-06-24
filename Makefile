.PHONY: service
service:
	go build -o service cmd/main.go  
	./service
.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.DEFAULT_GOAL := build