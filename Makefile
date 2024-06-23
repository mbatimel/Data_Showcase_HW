.PHONY: server
server:
	go build -o server cmd/main.go  
	./server
.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.DEFAULT_GOAL := build