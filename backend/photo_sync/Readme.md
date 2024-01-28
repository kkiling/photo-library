Download goose utils
```
go get -u github.com/pressly/goose/v3/cmd/goose
```

Create migration
```
goose create [name] sql
```
SQLite migration
```
goose -dir=./migrations sqlite3 ./files.db up
goose -dir=./migrations sqlite3 ./files.db down
```
