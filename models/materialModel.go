package models

import (
	"time"

	"gorm.io/gorm"
)

type Material struct {
	gorm.Model
	Title        string
	Poster       string
	Description  string
	Publish_year int
	Director     string
	Producer     string
	Categories   []*Category `gorm:"many2many:material_categories;"`
	Ages         []*Age      `gorm:"many2many:material_ages;"`
	Genres       []*Genre    `gorm:"many2many:material_genres;"`
	Image_src    []Image_src `gorm:"foreignKey:Material_id;references:ID"`
	Duration     string
	Keywords     string
	M_type       string
}

type Material_recommend struct {
	Material_id uint
	Title       string
	Poster      string
	Description string
}

type Movie struct {
	Id           uint `json:"id"`
	Poster       string
	Title        string
	Publish_year int
	Duration     string
	Description  string
	Director     string
	Producer     string
}

type Material_history struct {
	Id     uint
	Poster string
	Title  string
}

type Material_get struct {
	Material_id uint
	Title       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Poster      string
	Category    string
}

type Material_search struct {
	Material_id  uint
	Title        string
	Poster       string
	Category     string
	Publish_year int
}

type Materials_get struct {
	Material_id uint
	Title       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Poster      string
	Category    []string
	Viewed      int
	Sezon       int
	Series      int
}
