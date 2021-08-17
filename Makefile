postgres:
	docker run --name postgres12 -p 5555:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

postgresdrop:
	docker stop postgres12 && docker rm postgres12

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_wallet

dropdb:
	docker exec -it postgres12 dropdb simple_wallet

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5555/simple_wallet?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5555/simple_wallet?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

sleep:
	sleep 5s

reinit: postgresdrop postgres sleep createdb migrateup

init: postgres sleep createdb migrateup

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test