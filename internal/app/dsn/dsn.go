// Package dsn собирает строку подключения к PostgreSQL из переменных
// окружения.
package dsn

import (
	"fmt"
	"os"
)

// FromEnv возвращает DSN для подключения к PostgreSQL, собранный из
// переменных окружения DB_HOST/DB_PORT/DB_USER/DB_PASS/DB_NAME.
func FromEnv() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return ""
	}
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}
