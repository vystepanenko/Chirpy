include ./.env
export

migrate_up:
	goose -dir sql/schema postgres "${DB_URL}" up

migrate_down:
	goose -dir sql/schema postgres "${DB_URL}" down

db_login:
	docker exec -it gator-db-1 psql -U ${DB_USER}

generate_sql:
	sqlc generate

build_and_run:
	go build && ./Chirpy
