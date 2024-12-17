include .env
migrate-up:
	@migrate -database ${DB_URL} -path ./internal/infra/migrations up

migrate-down:
	@migrate -database ${DB_URL} -path ./internal/infra/migrations down

dev:
	@docker compose up -d --build && \
	sleep 5 && \
	echo 'Running Migration...' && \
	migrate -database ${DB_URL} -path ./internal/infra/migrations up && \
	air

clean:
	@docker compose down --remove-orphans --volumes
