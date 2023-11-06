createmigration:
	migrate create -ext=sql -dir=sql/migrations -seq init
migrate:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/course" -verbose up
migratedown:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/course" -verbose down

.PHONY: migrate migrate-down createmigration