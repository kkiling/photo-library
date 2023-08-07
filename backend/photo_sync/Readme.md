```
go get -u github.com/pressly/goose/v3/cmd/goose
```

Create migration
```
goose create [name] sql
```
```
goose -dir=./migrations sqlite3 ./files.db up
goose sqlite3 ./files.db down
```
