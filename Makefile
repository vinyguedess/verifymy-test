
dev:
	go run .

down:
	docker-compose stop
	docker-compose down

pre-test:
	mkdir -p coverage
	make pre-test-build

pre-test-build:
	rm -rf mocks

test:
	make pre-test
	2>&1 ENV=test go test -v ./... -v -coverprofile ./coverage/.coverage
	go tool cover -html=./coverage/.coverage -o ./coverage/index.html
	go tool cover -func ./coverage/.coverage

up:
	docker-compose up -d
	docker-compose exec app bash
