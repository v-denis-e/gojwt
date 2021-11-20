# Перед использованием make необходимо создать локальную версию
# .env файла. Для шаблона можно использовать .env.default
include .env

# Запускает приложение с исходников: make run
run:
	go run .

# Запускает все тесты: make test
test:
	go test ./... -v

# Собирает проект: make build
build:
	go build .