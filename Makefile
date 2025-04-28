default: run

.PHONY: run
run: # Execute the Go server
	@go run cmd/api/main.go

.PHONY: migrate
migrate: # Add a new migration
	@echo "====> Adding a new migration"
	@if [ -z "$(name)" ]; then echo "Migration name is required"; exit 1; fi
	@migrate create -ext sql -dir services/internal/infra/database/migrate/migrations $(name)

.PHONY: migrate-up
migrate-up: # Apply all pending migrations
	@echo "====> Applying all pending migrations"
	@go run cmd/migrate/main.go up

.PHONY: migrate-down
migrate-down: # Revert all applied migrations
	@echo "====> Reverting all applied migrations"
	@go run cmd/migrate/main.go down
