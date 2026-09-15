package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"utility-bills-backend/internal/app/ds"
)

// GetCommunalResources — список опубликованных услуг для страницы «Плитка»,
// отсортированный по ID. Если maxTariffRate указан, остаются только услуги
// с тарифом не выше этого значения (фильтр по теме).
func (r *Repository) GetCommunalResources(maxTariffRate *float64) ([]ds.CommunalResource, error) {
	query := r.db.Where("status = ?", ds.StatusPublished).Order("id")
	if maxTariffRate != nil {
		query = query.Where("tariff_rate <= ?", *maxTariffRate)
	}

	var resources []ds.CommunalResource
	if err := query.Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

// GetCommunalResourceFeedItem — услуга для страницы «Лента». Если next
// равен false, возвращается сама услуга по id (id == 0 — первая
// опубликованная, для перехода из панели вкладок). Если next равен true —
// следующая после неё опубликованная услуга, после последней — снова первая.
func (r *Repository) GetCommunalResourceFeedItem(id uint, next bool) (ds.CommunalResource, error) {
	var published []ds.CommunalResource
	if err := r.db.Where("status = ?", ds.StatusPublished).Order("id").Find(&published).Error; err != nil {
		return ds.CommunalResource{}, err
	}
	if len(published) == 0 {
		return ds.CommunalResource{}, gorm.ErrRecordNotFound
	}

	if id == 0 {
		return published[0], nil
	}

	index := -1
	for i, res := range published {
		if res.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return ds.CommunalResource{}, gorm.ErrRecordNotFound
	}

	if next {
		index = (index + 1) % len(published)
	}
	return published[index], nil
}

// GetCommunalResourceDraft — черновик текущего создателя (не более одного
// на пользователя). Отсутствие черновика — не ошибка, возвращается nil.
func (r *Repository) GetCommunalResourceDraft(creatorID uint) (*ds.CommunalResource, error) {
	var draft ds.CommunalResource
	err := r.db.Where("status = ? AND creator_id = ?", ds.StatusDraft, creatorID).First(&draft).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &draft, nil
}

// CreateCommunalResourceDraft создаёт черновик услуги по кнопке «Далее».
// Фото и видео в лабораторной 2 не сохраняются (требование методички) —
// принимается только название.
func (r *Repository) CreateCommunalResourceDraft(creatorID uint, name string) (*ds.CommunalResource, error) {
	draft := ds.CommunalResource{
		Name:      name,
		Status:    ds.StatusDraft,
		CreatorID: creatorID,
		CreatedAt: time.Now(),
	}
	if err := r.db.Create(&draft).Error; err != nil {
		return nil, err
	}
	return &draft, nil
}

// PublishCommunalResource заполняет оставшиеся поля черновика и переводит
// его в статус «опубликован» по кнопке «Опубликовать».
func (r *Repository) PublishCommunalResource(id uint, description string, tariffRate, consumptionNorm float64, measurementUnit string) error {
	now := time.Now()
	return r.db.Model(&ds.CommunalResource{}).
		Where("id = ? AND status = ?", id, ds.StatusDraft).
		Updates(map[string]any{
			"description":      description,
			"tariff_rate":      tariffRate,
			"consumption_norm": consumptionNorm,
			"measurement_unit": measurementUnit,
			"status":           ds.StatusPublished,
			"published_at":     now,
		}).Error
}

// DeleteCommunalResource — логическое удаление услуги прямым SQL UPDATE,
// без ORM (требование методички лабораторной 2).
func (r *Repository) DeleteCommunalResource(id uint) error {
	return r.db.Exec("UPDATE communal_resources SET status = ? WHERE id = ?", ds.StatusDeleted, id).Error
}

// GetLikesCount возвращает количество лайков услуги — вычисляется по
// таблице likes, а не хранится отдельным полем.
func (r *Repository) GetLikesCount(resourceID uint) (int64, error) {
	var count int64
	err := r.db.Model(&ds.Like{}).Where("communal_resource_id = ?", resourceID).Count(&count).Error
	return count, err
}
