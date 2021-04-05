# Documentation

Table of Contents
=================
 * [File Structure](#File-Structure)
 * [Files](#Files)
 * [Migration](#Migration)
 * [Util](#Util)

## File Structure
```
project
└── api
    └── contact.go
└── migrations
    ├── 000001_create_bike_table.down.sql
    └── 000001_create_bike_table.up.sql
└── models
    └── error.go
└── postgres
    └── rows.go
└── templates
    └── register.html
└── util
    ├── log.go
    └── ...
└── applicaiton.go
└── config.json
└── docker-compose.yml
└── go.mod
└── go.sum
```

## Files
1. api 
  ```
   Contains all api handles needed for the project. Make sure one(1) file per use case.
   Ex. user.go, contact.go, orders.go, etc.
   ```
2. migrations
   ```
   Contains the sql files craeted with the go-migration CLI.
   ```
3. models
   ```
   Contains all global structs, place here all structs that are needed in more than one go file so we can re use it. Also separate files based on the use case.
   Ex. DefaultError is used in different go files.
   ```
4. postgres
   ```
   Contains generic functions for postgres query. If you use the same methods and queries for multiple postgres calls, better place it here.
   ```
5. templates
   ```
   Contains html templates and other stuffs. This template can be used during  sending emails and etc.
   ```
6. util
   ```
   Contains utility functions and helpers that you will use.
   Example of this is validation, error handles, writing responses, creating context
   ```
7. application.go
   ```
   This is our main.go, basically you will only modify the route Handle Functions here. Declare all endpoints and create a new context pointing to the go file in the api folder.
   ```
8. config.json
   ```
   This contains senstive informations like keys, postgres connection string, secrets and etc.
   Basically the .env file
   ```
9. docker-compose.yml
   ```
   Docker compose configuration for the postgres
   ```
10. go.mod
11. go.sum

## Migration
To know more about go migration visit this page: [Go-Migrate](https://github.com/golang-migrate/migrate)

##### Create some migrations using migrate CLI. Here is an example:
```
migrate create -ext sql -dir ./migrations -seq create_contacts_table
```

##### Run your migrations through the CLI or your app and check if they applied expected changes. Just to give you an idea:
```
migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -source file://migrations up
```

##### Forcing your database version
> In case you run a migration that contained an error, migrate will not let you run other migrations on the same database. You will see an error like Dirty database version 1. Fix and force version, even when you fix the erred migration. This means your database was marked as 'dirty'. You need to investigate the migration error - was your migration applied partially, or was it not applied at all? Once you know, you should force your database to a version reflecting it's real state. You can do so with force command:
For more details check this [comment](https://github.com/golang-migrate/migrate/issues/282#issuecomment-530743258)
```
error: Dirty database version 16. Fix and force version.


migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -source file://migrations force VERSION

## VERION -> replace with the latest working version, in our example it is version 15, because version 16 had an issue.

migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -source file://migrations force 15
```

## Util
(WIP)
