// Точка входа приложения. По методическим указаниям курса файл main.go
// должен содержать минимум логики — весь запуск сервера вынесен в
// internal/api.
package main

import (
	"log"

	"utility-bills-backend/internal/api"
)

func main() {
	log.Println("Application start!")
	api.StartServer()
	log.Println("Application terminated!")
}
