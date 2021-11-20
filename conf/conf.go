// conf пакет для получения конфигурации приложения

package conf

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

const (
	postgresHost        = "POSTGRES_HOST"
	defaultPostgresHost = "localhost"

	postgresPort        = "POSTGRES_PORT"
	defaultPostgresPort = "5432"

	postgresDB        = "POSTGRES_DB"
	defaultPostgresDB = "postgres"

	postgresUser        = "POSTGRES_USER"
	defaultPostgresUser = "postgres"

	postgresPassword = "POSTGRES_PASSWORD"
)

// AppConf представляет главную конфигурацию приложения
type AppConf struct {
	DB DBConf
}

// DBConf представляет конфигурацию для подключения к бд
type DBConf struct {
	Host     string // Хост сервера БД
	Port     int    // Порт подключения к БД
	Name     string // Имя базы данных
	User     string // Имя пользователя
	Password string // Пароль от пользователя
}

// FromEnv eизвлекает конфигурацию приложения из переменных среды
func FromEnv() (AppConf, error) {
	log.Info().Msg("Getting app configuration from environment ...")

	db, err := extractDbConf()
	if err != nil {
		return AppConf{}, err
	}

	return AppConf{DB: db}, nil
}

// extractDbConf извлекает параметры подключения к бд
func extractDbConf() (DBConf, error) {
	host := extractOrDefault(postgresHost, defaultPostgresHost)

	portParam := extractOrDefault(postgresPort, defaultPostgresPort)
	port, err := strconv.Atoi(portParam)
	if err != nil {
		return DBConf{}, err
	}

	name := extractOrDefault(postgresDB, defaultPostgresDB)

	user := extractOrDefault(postgresUser, defaultPostgresUser)

	password := os.Getenv(postgresPassword)

	return DBConf{host, port, name, user, password}, nil
}

// extractOrDefault извлекает по ключу значение из среды или возвращает значение по умолчанию
func extractOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		debug(key, defaultValue)
		return defaultValue
	}

	return value
}

func debug(key string, value string) {
	log.Debug().Msg(fmt.Sprintf("%s is empty setting default value: %s", key, value))
}
