default:
	docker-compose build auth

up: default
	docker-compose up

down:
	docker-compose down

test:
	go test -v -cover ./...

clean: down
	rm -f hot_reload_exec
	docker system prune -f
	docker volume prune -f

migrateup:
	docker-compose run --rm migrator "make up"

migratedown:
	docker-compose run --rm migrator "make down"

.PHONY: up down test clean migrateup migratedown
