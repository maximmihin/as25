# Сделано
## Ручки (с e2e тестами)
- [x] post /dummyLogin
- [x] post /pvz
- [x] post /receptions
- [x] post /products
- [x] post /pvz/{pvzId}/close_last_reception
- [x] post /pvz/{pvzId}/delete_last_product
- [ ] get /pvz (50%)
- [ ] post /register
- [ ] post /login

## Дополнительные задания
- [x] логирование [вдохновлено](https://www.youtube.com/watch?v=p9XiOOU52Qw)
- [x] кодогенерация моделей (+бойлерплейта)
- [ ] /register и /login
- [ ] prometheus
- [ ] gRPC
-------------------------------------------------------------------------------
- [ ] поправлена спека (schemas, requests, enums etc)


## Структура проекта
```shell
.
├── README.md
├── api
│   └── api.yaml
├── cmd
│   └── httpserver
│       ├── e2e_test.go
│       ├── main.go
│       └── testclient            // обертка для сгенерированного клиента для e2e тестов
│           ├── cfg.yaml
│           ├── client.gen.go
│           └── testclient.go
├── go.mod
├── go.sum
├── internal
│   ├── controllers
│   │   └── http
│   │       ├── handlers.go
│   │       ├── middleware
│   │       │   ├── authorization.go
│   │       │   └── error_handler.go
│   │       ├── models_and_boilerplate.gen.go
│   │       ├── oapi_config.yaml
│   │       ├── requests_defaults.go
│   │       └── requests_validators.go
│   ├── dal
│   │   ├── embed_migrations.go        // emped для включения sql файлов как файловой системы для запуска миграций из кода
│   │   ├── migrations
│   │   │   └── 00001_init.sql
│   │   ├── queries                   // queries (недописанный) - для разделения command & queries. queries - для сложных запросов на извлечение
│   │   │   ├── models
│   │   │   │   ├── db.go
│   │   │   │   ├── models.go
│   │   │   │   └── queries.sql.go
│   │   │   ├── queries.go
│   │   │   ├── queries.sql
│   │   │   └── sqlc.yaml
│   │   ├── repos
│   │   │   ├── products
│   │   │   │   ├── models
│   │   │   │   │   ├── db.go
│   │   │   │   │   ├── models.go
│   │   │   │   │   └── queries.sql.go
│   │   │   │   ├── queries.sql
│   │   │   │   ├── repo.go
│   │   │   │   └── sqlc.yaml
│   │   │   ├── pvz
│   │   │   │   ├── models
│   │   │   │   │   ├── db.go
│   │   │   │   │   ├── models.go
│   │   │   │   │   └── queries.sql.go
│   │   │   │   ├── queries.sql
│   │   │   │   ├── repo.go
│   │   │   │   └── sqlc.yaml
│   │   │   └── receptions
│   │   │       ├── models
│   │   │       │   ├── db.go
│   │   │       │   ├── models.go
│   │   │       │   └── queries.sql.go
│   │   │       ├── queries.sql
│   │   │       ├── repo.go
│   │   │       └── sqlc.yaml
│   │   └── types
│   │       └── city.go
│   ├── logger
│   │   └── ctx_enricher.go
│   └── services
│       ├── accesscontrol
│       │   ├── dummy_jwt.go
│       │   └── pkg
│       │       └── jwt
│       │           └── custom_claims.go
│       └── pvz
└── task.md
```
