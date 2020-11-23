package service

import "card-keeper-api/model"

// Service for interacting with data store
type Service struct {
	Repository
}

// GetAll returns all Cards in the repository
func (service *Service) GetAll() *[]model.Card {
	cards, _ := service.Repository.GetAll()
	return cards
}
