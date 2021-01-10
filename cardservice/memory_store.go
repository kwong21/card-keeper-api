package cardservice

import (
	"reflect"
)

type memoryStore struct {
	Cards []Card
}

//InMemoryStore returns an in-memory repository
func InMemoryStore() (Repository, error) {
	return &memoryStore{
		Cards: make([]Card, 0),
	}, nil
}

func (r *memoryStore) GetAll() (*[]Card, error) {
	return &r.Cards, nil
}

func (r *memoryStore) AddCard(card Card) error {
	var err error

	for _, c := range r.Cards {
		if card.Base.Player == c.Base.Player {
			if reflect.DeepEqual(card, c) {
				err = &DuplicateError{}
			}
		}
	}

	if err == nil {
		r.Cards = append(r.Cards, card)
	}

	return err
}
