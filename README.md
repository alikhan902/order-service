# Order Service

[![CI](https://github.com/alikhan902/order-service/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/alikhan902/order-service/actions/workflows/ci.yml)
[![Docker](https://github.com/alikhan902/order-service/actions/workflows/docker.yml/badge.svg?branch=main)](https://github.com/alikhan902/order-service/actions/workflows/docker.yml)

Сервис принимает заказы по HTTP на порту 8080 и хранит их в памяти.

## Проверка

    go build ./...
    go test -race -coverpkg=./... -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

CI выполняет эти проверки и golangci-lint для Go 1.21 и 1.22. Профиль покрытия каждой версии доступен в разделе Artifacts запуска CI.

## Контейнер

    docker build -t order-service:lab7 .
    docker run --rm -p 8080:8080 order-service:lab7

После успешной проверки push в main workflow Docker публикует теги main и sha-<первые 7 символов коммита>. Для публикации нужны секреты репозитория DOCKERHUB_USERNAME и DOCKERHUB_TOKEN с правом записи в Docker Hub.
