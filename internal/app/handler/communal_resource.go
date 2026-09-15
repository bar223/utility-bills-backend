package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"utility-bills-backend/internal/app/ds"
)

// defaultImageURL — изображение по умолчанию, хранится вместе с иконками
// на SSR-сервере (не в MinIO) и подставляется, если у услуги нет своего
// изображения или указанный файл недоступен.
const defaultImageURL = "/static/img/default-image.svg"

// resourceView — данные услуги для шаблонов: сама запись плюс количество
// лайков (вычисляется по таблице likes) и изображение с фолбэком на
// значение по умолчанию.
type resourceView struct {
	ds.CommunalResource
	LikesCount      int64
	DisplayImageURL string
}

func (h *Handler) toView(r ds.CommunalResource) resourceView {
	likes, err := h.Repository.GetLikesCount(r.ID)
	if err != nil {
		logrus.Error(err)
	}

	imageURL := r.ImageURL
	if imageURL == "" {
		imageURL = defaultImageURL
	}

	return resourceView{
		CommunalResource: r,
		LikesCount:       likes,
		DisplayImageURL:  imageURL,
	}
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
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	views := make([]resourceView, 0, len(resources))
	for _, r := range resources {
		views = append(views, h.toView(r))
	}

	ctx.HTML(http.StatusOK, "tile.html", gin.H{
		"Resources": views,
		"Tariff":    tariffQuery,
	})
}

// Feed — GET /feed и GET /feed/:id: страница «Лента», один
// CommunalResource на экран, ?next=true — переход к следующему.
func (h *Handler) Feed(ctx *gin.Context) {
	idParam := ctx.Param("id")

	var id uint
	if idParam != "" {
		parsed, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			ctx.String(http.StatusBadRequest, "некорректный идентификатор услуги")
			return
		}
		id = uint(parsed)
	}

	next := ctx.Query("next") == "true"

	resource, err := h.Repository.GetCommunalResourceFeedItem(id, next)
	if err != nil {
		ctx.String(http.StatusNotFound, "услуга не найдена")
		return
	}

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"Resource": h.toView(resource),
	})
}

// Add — GET /add: черновик текущего создателя. Если черновика ещё нет,
// показывается форма создания (кнопка «Далее»), если уже есть — форма
// публикации с заполненными полями (кнопка «Опубликовать»).
func (h *Handler) Add(ctx *gin.Context) {
	draft, err := h.Repository.GetCommunalResourceDraft(currentCreatorID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "add.html", gin.H{
		"Draft": draft,
	})
}

// CreateDraft — POST /add, кнопка «Далее»: создаёт черновик по названию.
// Фото и видео в лабораторной 2 не сохраняются (требование методички).
func (h *Handler) CreateDraft(ctx *gin.Context) {
	name := ctx.PostForm("name")

	_, err := h.Repository.CreateCommunalResourceDraft(currentCreatorID, name)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/add")
}

// Publish — POST /publish, кнопка «Опубликовать»: заполняет оставшиеся
// поля черновика и меняет статус на «опубликован».
func (h *Handler) Publish(ctx *gin.Context) {
	idStr := ctx.PostForm("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	description := ctx.PostForm("description")
	measurementUnit := ctx.PostForm("measurementUnit")

	tariffRate, err := strconv.ParseFloat(ctx.PostForm("tariffRate"), 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	consumptionNorm, err := strconv.ParseFloat(ctx.PostForm("consumptionNorm"), 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.PublishCommunalResource(uint(id), description, tariffRate, consumptionNorm, measurementUnit); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}

// Delete — POST /delete-resource: логическое удаление услуги с плитки
// прямым SQL UPDATE.
func (h *Handler) Delete(ctx *gin.Context) {
	idStr := ctx.PostForm("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.DeleteCommunalResource(uint(id)); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/")
}
