
dev:
	go run .

down:
	docker-compose stop
	docker-compose down

up:
	docker-compose up -d
	docker-compose exec app bash
