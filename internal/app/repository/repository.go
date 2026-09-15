// Package repository отвечает за доступ к данным. В первой лабораторной
// работе это не БД, а коллекция в памяти — требование методических
// указаний курса (лаба 1: "все данные для страниц нужно брать прямо из
// коллекции ... без использования БД").
package repository

import (
	"errors"
	"time"
)

// Статусы CommunalResource — жизненный цикл записи справочника услуг.
const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusDeleted   = "deleted"
)

// CommunalResource — коммунальный ресурс (услуга), основная сущность
// варианта «Квитанция по квартплате».
type CommunalResource struct {
	ID              int
	Name            string
	TariffRate      float64
	MeasurementUnit string
	Description     string
	ImageURL        string
	VideoURL        string
	Likes           []int // ID пользователей, поставивших лайк
	Status          string
	CreatedAt       time.Time
}

var (
	// ErrCommunalResourceNotFound возвращается, когда запись по ID не найдена
	// среди опубликованных ресурсов.
	ErrCommunalResourceNotFound = errors.New("услуга не найдена")
	// ErrDraftNotFound возвращается, когда в коллекции нет черновика.
	ErrDraftNotFound = errors.New("черновик не найден")
)

const minioBaseURL = "http://localhost:9000/utility-bills/"

// Repository хранит коллекцию CommunalResource в памяти. В следующих
// лабораторных работах здесь появится подключение к PostgreSQL, сигнатуры
// методов при этом не изменятся.
type Repository struct {
	communalResources []CommunalResource
}

// NewRepository создаёт репозиторий и наполняет его тестовыми данными по
// шести коммунальным ресурсам варианта плюс одним черновиком.
func NewRepository() (*Repository, error) {
	r := &Repository{
		communalResources: []CommunalResource{
			{
				ID:              1,
				Name:            "Электроэнергия",
				TariffRate:      5.68,
				MeasurementUnit: "кВт·ч",
				Description:     "Оплата электрической энергии, потребляемой квартирой: освещение, розетки и бытовые приборы.",
				ImageURL:        minioBaseURL + "electricity.webp",
				VideoURL:        minioBaseURL + "electricity.mp4",
				Likes:           []int{101, 102, 103, 104, 105},
				Status:          StatusPublished,
				CreatedAt:       time.Date(2026, 1, 10, 9, 0, 0, 0, time.UTC),
			},
			{
				ID:              2,
				Name:            "Холодное водоснабжение",
				TariffRate:      47.24,
				MeasurementUnit: "м³",
				Description:     "Оплата холодной воды, поступающей в кран и смывной бачок, по показаниям счётчика.",
				ImageURL:        minioBaseURL + "cold-water.avif",
				VideoURL:        minioBaseURL + "cold-water.mp4",
				Likes:           []int{101, 102, 103, 104},
				Status:          StatusPublished,
				CreatedAt:       time.Date(2026, 1, 10, 9, 5, 0, 0, time.UTC),
			},
			{
				ID:              3,
				Name:            "Капитальный ремонт",
				TariffRate:      24.2,
				MeasurementUnit: "м²",
				Description:     "Взнос на капитальный ремонт общего имущества многоквартирного дома.",
				ImageURL:        minioBaseURL + "capital-repair.webp",
				VideoURL:        minioBaseURL + "capital-repair.mp4",
				Likes:           []int{101, 102, 103, 104, 105, 106, 107, 108},
				Status:          StatusPublished,
				CreatedAt:       time.Date(2026, 1, 10, 9, 10, 0, 0, time.UTC),
			},
			{
				ID:              4,
				Name:            "Вывоз ТКО",
				TariffRate:      6.5,
				MeasurementUnit: "м²",
				Description:     "Оплата вывоза твёрдых коммунальных отходов управляющей организацией.",
				ImageURL:        minioBaseURL + "waste-removal.jpg",
				VideoURL:        minioBaseURL + "waste-removal.mp4",
				Likes:           []int{101, 102, 103, 104, 105, 106, 107, 108, 109, 110},
				Status:          StatusPublished,
				CreatedAt:       time.Date(2026, 1, 10, 9, 15, 0, 0, time.UTC),
			},
			{
				ID:              5,
				Name:            "Газоснабжение",
				TariffRate:      8.2,
				MeasurementUnit: "м³",
				Description:     "Оплата природного газа, который используется для плиты и проточного водонагревателя.",
				ImageURL:        minioBaseURL + "gas-supply.avif",
				VideoURL:        minioBaseURL + "gas-supply.mp4",
				Likes:           []int{101, 102, 103, 104, 105, 106},
				Status:          StatusPublished,
				CreatedAt:       time.Date(2026, 1, 10, 9, 20, 0, 0, time.UTC),
			},
			{
				ID:              6,
				Name:            "Отопление",
				TariffRate:      2450,
				MeasurementUnit: "Гкал",
				Description:     "Начисление за отопление помещения рассчитывается по нормативу потребления на квадратный метр площади и действует весь отопительный сезон, включая летние месяцы по решению УК.",
				ImageURL:        minioBaseURL + "heating.webp",
				VideoURL:        minioBaseURL + "heating.mp4",
				Likes:           []int{101, 102, 103, 104, 105},
				Status:          StatusPublished,
				CreatedAt:       time.Date(2026, 1, 10, 9, 25, 0, 0, time.UTC),
			},
			{
				// Черновик: создатель начал заполнять новую услугу "Отопление",
				// фото и видео ещё не приложены — на странице показываются
				// иконки-заглушки. Записи в статусе draft не показываются
				// ни в плитке, ни в ленте.
				ID:              7,
				Name:            "Отопление",
				TariffRate:      2450,
				MeasurementUnit: "Гкал",
				Description:     "Начисление за отопление помещения рассчитывается по нормативу потребления на квадратный метр площади и действует весь отопительный сезон, включая летние месяцы по решению УК.",
				ImageURL:        "",
				VideoURL:        "",
				Likes:           nil,
				Status:          StatusDraft,
				CreatedAt:       time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC),
			},
		},
	}
	return r, nil
}

// published возвращает опубликованные ресурсы в порядке возрастания ID —
// в этом же порядке они идут в плитке и в ленте.
func (r *Repository) published() []CommunalResource {
	result := make([]CommunalResource, 0, len(r.communalResources))
	for _, cr := range r.communalResources {
		if cr.Status == StatusPublished {
			result = append(result, cr)
		}
	}
	return result
}

// GetCommunalResources возвращает список опубликованных услуг для страницы
// «Плитка». Если maxTariffRate указан, остаются только услуги с тарифом не
// выше этого значения — единственное числовое поле по теме, доступное для
// фильтрации на сервере.
func (r *Repository) GetCommunalResources(maxTariffRate *float64) ([]CommunalResource, error) {
	all := r.published()
	if maxTariffRate == nil {
		return all, nil
	}

	filtered := make([]CommunalResource, 0, len(all))
	for _, cr := range all {
		if cr.TariffRate <= *maxTariffRate {
			filtered = append(filtered, cr)
		}
	}
	return filtered, nil
}

// GetCommunalResourceFeedItem возвращает услугу для страницы «Лента».
// Если next=false — саму услугу с указанным ID. Если next=true — следующую
// после неё опубликованную услугу; после последней лента зацикливается на
// первую. Если id равен нулю (переход из панели вкладок без указания ID),
// возвращается первая опубликованная услуга.
func (r *Repository) GetCommunalResourceFeedItem(id int, next bool) (CommunalResource, error) {
	published := r.published()
	if len(published) == 0 {
		return CommunalResource{}, ErrCommunalResourceNotFound
	}

	if id == 0 {
		return published[0], nil
	}

	index := -1
	for i, cr := range published {
		if cr.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		return CommunalResource{}, ErrCommunalResourceNotFound
	}

	if next {
		index = (index + 1) % len(published)
	}
	return published[index], nil
}

// GetCommunalResourceDraft возвращает единственную услугу в статусе
// «черновик» — она же отображается на странице «Добавление».
func (r *Repository) GetCommunalResourceDraft() (CommunalResource, error) {
	for _, cr := range r.communalResources {
		if cr.Status == StatusDraft {
			return cr, nil
		}
	}
	return CommunalResource{}, ErrDraftNotFound
}
