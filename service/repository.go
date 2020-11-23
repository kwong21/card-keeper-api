package service

import "card-keeper-api/model"

// Repository writes and reads data from declared store
type Repository interface {
	GetAll() (*[]model.Card, error)
	AddCard(model.Card) error
}
