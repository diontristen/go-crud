# GO-TEMPLATE
This repository can be forked/clone when starting a new GO Project.

## Before Starting to develop:

1. Rename the parent directory
```
  Command: mv go-template <project_name>
  Example: mv go-tempalte test-backend
```
 

2. Delete .git file
```
  Command :    rm -rf .git
  Description: This will delete the .git file of this template
```

3. Initialize Git
```
  Command:    git init
  Desription: To initialize your new git
```

4. Create .gitignore
```
  Command: touch .gitignore
  Copy ./gitignore.txt
  ## Make sure to ignore sensitivie files
```

5. Rename (manually) all the local import package declaration to base on your project (looking for a better way or to automate this)
```
FROM: github.com/diontristen/go-template/*
TO:   github.com/<username>/<repositroy>/*
EX:   github.com/team-yl/test/*

FILES TO UPDATE:
 - application.go
 - api/contact.go
 - postgres/rows.go
 - util/
   - error.go
   - errors.go
   -  
```
6. Remove unnecessary files
```
Ex: api/contact.go
    migrations/*
``` 

## Documentation
If this is  your first time to create a go project. Make sure to check out the [documentation](documentation.md)

## Sample Project:
[Contact-CRUD-API](github.com/diontristen/contact-go)

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

Finally you can start the server:
```
go run application.go
```

   

