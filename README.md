# Каталог автомобилей

## Конфигурация

```
mv .env.example .env
mv config/example.toml config/config.toml
```

### Внешний API

Ссылку на внешний API необходимо указывать по пути `config/config.toml` в `api.url`.

## Запуск

### Сервис

```
go run ./cmd/main.go
```

### Документация

```
http://localhost:8081/api/v1/swagger/index.html
```

### Docker

```
docker-compose up -d
```

## Миграции

Миграции автоматически создаются после запуска Docker.
