# Order Service Backend

## How to Run

go run cmd/api/main.go

## How to Migrate the Database

migrate -path database/migrations/ -database "mysql://DB_USERNAME:DB_PASSWORD@tcp(localhost:DB_PORT)/DB_NAME" -verbose up
