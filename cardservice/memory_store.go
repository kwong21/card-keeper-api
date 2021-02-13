package cardservice

import (
	"reflect"
)

type memoryStore struct {
	Cards map[string][]Card
}

//InMemoryStore returns an in-memory repository
func InMemoryStore() (Repository, error) {
	return &memoryStore{
		Cards: make(map[string][]Card),
	}, nil
}

func (r *memoryStore) GetAllCardsInCollection(collection string) ([]Card, error) {
	return r.Cards[collection], nil
}

func (r *memoryStore) AddCardToCollection(card Card, collection string) error {
	var err error

	for _, c := range r.Cards[collection] {
		if card.Base.Player == c.Base.Player {
			if reflect.DeepEqual(card, c) {
				err = &DuplicateError{}
			}
		}
	}

	if err == nil {
		cardsInCollection := r.Cards[collection]
		r.Cards[collection] = append(cardsInCollection, card)
	}

	return err
}
