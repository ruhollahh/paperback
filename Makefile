## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## db/psql: connect to the database using psql
db/psql:
	psql ${PAPERBACK_DB_DSN}

## db/migrations/new name=$1: create a new database migration
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${PAPERBACK_DB_DSN} up

## db/migrations/down: rollback all database migrations
db/migrations/down:
	@echo 'Rollingback all migrations...'
	migrate -path ./migrations -database ${PAPERBACK_DB_DSN} down

## db/migrations/force version=$1: force a dirty version
db/migrations/force:
	@echo 'Forcing migration files for version ${version}...'
	migrate -path ./migrations -database ${PAPERBACK_DB_DSN} force ${version}