run:
	go run ./cmd/api

migrate:
	psql "$(POSTGRES_DSN)" -f db/migrations/000001_init.sql

sqlboiler:
	sqlboiler psql --config sqlboiler.toml