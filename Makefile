include .env
export

export PROJECT_ROOT=$(shell pwd)
env-up:
	docker compose up -d todoapp-postgres
env-down:
	docker compose down todoapp-postgres

env-cleanup:
	@read -p "Are you sure you want to delete the Postgres data? (y/n): " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Postgres data deleted."; \
	else \
		echo "Postgres data not deleted."; \
	fi

env-port-forwarder:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Error: Please provide a migration name using the 'seq' variable."; \
		exit 1; \
	fi

	docker compose run --rm todoapp-postgres-migrate \
	 create \
	 -ext sql \
	 -dir /migrations \
	 -seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Error: Please provide a migration action using the 'action' variable."; \
		exit 1; \
	fi
	docker compose run --rm todoapp-postgres-migrate \
	 -path /migrations \
	 -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable \
	 "$(action)"
	 
todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go