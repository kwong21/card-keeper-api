package service

import (
	"card-keeper-api/model"

	"gorm.io/gorm"
)

type dbStore struct {
	Db *gorm.DB
}

// Database provides access to the backend configured for webapp
func Database(db *gorm.DB) Repository {
	return &dbStore{
		Db: db,
	}
}

func (r *dbStore) GetAll() (*[]model.Card, error) {
	cards := []model.Card{}
	r.Db.Find(&cards)

	return &cards, nil
}

func (r *dbStore) AddCard(card model.Card) error {
	dbc := r.Db.Create(&card)

	return dbc.Error
}
