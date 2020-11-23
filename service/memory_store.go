package service

import "card-keeper-api/model"

type memoryStore struct {
	Cards []model.Card
}

//InMemoryStore returns an in-memory repository
func InMemoryStore() Repository {
	return &memoryStore{
		Cards: make([]model.Card, 0),
	}
}

func (r *memoryStore) GetAll() (*[]model.Card, error) {
	return &r.Cards, nil
}

func (r *memoryStore) AddCard(card model.Card) error {
	r.Cards = append(r.Cards, card)

	return nil
}
