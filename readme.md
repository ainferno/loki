# Loki

## Start mysql service

```
docker compose up -d mysql
```

## Start go app locally

```
go run main.go
```

## Start go app inside docker container

1. change `DB_HOST` variable to mysql service

```
# .env
DB_HOST=mysql
```

2. start application container

```
docker compose up -d app
```
