package cardservice

import "errors"

type memoryStore struct {
	Cards map[string]map[uint32]Card
}

//InMemoryStore returns an in-memory repository
func InMemoryStore() (Repository, error) {
	return &memoryStore{
		Cards: make(map[string]map[uint32]Card),
	}, nil
}

func (r *memoryStore) GetCardInCollection(cardID uint32, collection string) (Card, error) {
	c, ok := r.Cards[collection]

	if !ok {
		return Card{}, errors.New("Collection does not exist")
	}

	card := c[cardID]

	return card, nil
}

func (r *memoryStore) GetWatchListedCards(collection string) ([]Card, error) {
	c, err := r.GetAllCardsInCollection(collection)
	watchListedCards := make([]Card, 0)

	if err != nil {
		return nil, err
	}

	for _, card := range c {
		if card.IsOnWatchList {
			watchListedCards = append(watchListedCards, card)
		}
	}

	return watchListedCards, err
}

func (r *memoryStore) GetAllCardsInCollection(collection string) ([]Card, error) {
	c := r.Cards[collection]
	v := make([]Card, 0, len(c))

	for _, card := range c {
		v = append(v, card)
	}
	return v, nil
}

func (r *memoryStore) AddCardToCollection(card Card, collection string) error {
	var err error

	if _, ok := r.Cards[collection]; !ok {
		cardsMap := make(map[uint32]Card)

		r.Cards[collection] = cardsMap

	}

	c := r.Cards[collection]

	if _, ok := c[card.CardID]; ok {
		err = &DuplicateError{}
	} else {
		c[card.CardID] = card
	}

	return err
}

func (r *memoryStore) UpdateCardInCollection(card Card, collection string) error {
	return nil
}
