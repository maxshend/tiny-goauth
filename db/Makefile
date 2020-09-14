up:
	migrate -path /migrations -database $(DB_URL) -verbose up

down:
	migrate -path /migrations -database $(DB_URL) -verbose down

.PHONY: up down
