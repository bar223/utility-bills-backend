// Package api собирает веб-сервис: маршруты, шаблоны и статику.
package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"utility-bills-backend/internal/app/handler"
	"utility-bills-backend/internal/app/repository"
)

// StartServer поднимает веб-сервис на порту 8080.
func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	h := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// Три GET-метода домена CommunalResource: плитка, лента, черновик.
	r.GET("/", h.Tile)
	r.GET("/feed", h.Feed)
	r.GET("/feed/:id", h.Feed)
	r.GET("/add", h.Add)

	if err := r.Run(":8080"); err != nil {
		logrus.Error(err)
	}

	log.Println("Server down")
}
