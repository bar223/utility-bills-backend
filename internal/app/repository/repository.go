// Package repository отвечает за доступ к данным — подключение к
// PostgreSQL через GORM.
package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Repository хранит подключение к БД.
type Repository struct {
	db *gorm.DB
}

// New открывает подключение к PostgreSQL по строке dsn.
func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}
