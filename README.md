## Initialization
First you need to start the database:
```
docker-compose up
```

After that you can install the dependencies
```
go get ./...
```

Then you need to migrate the database (do this if you have already created a migration):
```
migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -source file://migrations up
```

Create a config.json file

```
{
        "postgres_conn": "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
}
```

Finally you can start the server:
```
go run application.go
```
