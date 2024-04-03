up-service:
	@docker-compose up --build service-chat

down-service:
	@docker-compose down

postgres-migrate-up:
	@migrate -path ./internal/db/schema -database 'postgres://postgres:qwerty@localhost:5436/db_chat?sslmode=disable' up

postgres-migrate-down:
	@migrate -path ./internal/db/schema -database 'postgres://postgres:qwerty@localhost:5436/db_chat?sslmode=disable' down