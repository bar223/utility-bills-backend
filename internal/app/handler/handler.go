// Package handler содержит обработчики HTTP-запросов домена
// CommunalResource: получают параметры запроса, вызывают репозиторий и
// передают данные в html/template. Бизнес-логики и обращений к
// хранилищу здесь нет — только оркестрация.
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"utility-bills-backend/internal/app/repository"
)

// Handler объединяет обработчики страниц вокруг общего репозитория.
type Handler struct {
	Repository *repository.Repository
}

// NewHandler создаёт обработчик с внедрённым репозиторием.
func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

// Tile — GET /: страница «Плитка», список опубликованных
// CommunalResource в две колонки с фильтром по тарифу.
func (h *Handler) Tile(ctx *gin.Context) {
	tariffQuery := ctx.Query("tariff")

	var maxTariffRate *float64
	if tariffQuery != "" {
		if parsed, err := strconv.ParseFloat(tariffQuery, 64); err == nil {
			maxTariffRate = &parsed
		} else {
			logrus.Warnf("некорректное значение фильтра tariff=%q: %v", tariffQuery, err)
		}
	}

	resources, err := h.Repository.GetCommunalResources(maxTariffRate)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "tile.html", gin.H{
		"Resources": resources,
		"Tariff":    tariffQuery,
	})
}

// Feed — GET /feed и GET /feed/:id: страница «Лента», один
// CommunalResource на экран. Параметр next=true открывает следующую
// опубликованную услугу после указанного ID (с переходом на первую после
// последней). Без ID (переход из панели вкладок) открывается первая
// опубликованная услуга.
func (h *Handler) Feed(ctx *gin.Context) {
	idParam := ctx.Param("id")

	var id int
	if idParam != "" {
		parsed, err := strconv.Atoi(idParam)
		if err != nil {
			logrus.Error(err)
			ctx.String(http.StatusBadRequest, "некорректный идентификатор услуги")
			return
		}
		id = parsed
	}

	next := ctx.Query("next") == "true"

	resource, err := h.Repository.GetCommunalResourceFeedItem(id, next)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusNotFound, "услуга не найдена")
		return
	}

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"Resource": resource,
	})
}

// Add — GET /add: страница «Добавление», отображает единственный
// черновик CommunalResource. Сохранение в первой лабораторной работе не
// реализуется — поля можно заполнять, но без отправки на сервер.
func (h *Handler) Add(ctx *gin.Context) {
	draft, err := h.Repository.GetCommunalResourceDraft()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "add.html", gin.H{
		"Draft": draft,
	})
}
