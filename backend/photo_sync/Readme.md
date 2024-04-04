Download goose utils
```
go get -u github.com/pressly/goose/v3/cmd/goose
```

Sqlite error: `Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub`
```
go env -w CGO_ENABLED=1
Windows: choco install mingw -y
Ubuntu: apt-get install build-essential
```
Create migration
```
goose create [name] sql
```
SQLite migration
```
goose -dir=migrations sqlite3 ./files.db up
```
