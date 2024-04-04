Run `go mod tidy` to resolve the versions. Install by running

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

```
DELETE FROM coeffs_similar_photos;
DELETE FROM exif_photo_data;
DELETE FROM photo_groups_photos;
DELETE FROM photo_groups;
DELETE FROM photo_metadata;
DELETE FROM photo_previews;
DELETE FROM photo_vectors;
DELETE FROM rocket_locks;
DELETE FROM tags;
DELETE FROM tags_category;
DELETE FROM photo_processing_statuses;
```