Run `go mod tidy` to resolve the versions. Install by running
```
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

Пример метрик
```
	// Создайте метрику
	myCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "my_custom_counter",
		Help: "This is my counter",
	})

	// Регистрируйте метрику
	prometheus.MustRegister(myCounter)

	// Обновите значение метрики
	myCounter.Inc()
```


```
go get -u github.com/pressly/goose/v3/cmd/goose
```

Create migration
```
goose create [name] sql
```
```
goose -dir=./migrations postgres "postgresql://localhost:5432/photo_library?sslmode=disable" up
goose-dir=./migrations postgres "postgresql://localhost:5432/photo_library?sslmode=disable" down
```
