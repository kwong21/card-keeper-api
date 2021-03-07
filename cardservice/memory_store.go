package cardservice

type memoryStore struct {
	Cards map[string]map[uint32]Card
}

//InMemoryStore returns an in-memory repository
func InMemoryStore() (Repository, error) {
	return &memoryStore{
		Cards: make(map[string]map[uint32]Card),
	}, nil
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
