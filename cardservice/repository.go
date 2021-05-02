package cardservice

// Service for interacting with data store
type Service struct {
	Repository
}

// Repository writes and reads data from declared store
type Repository interface {
	GetCardInCollection(uint32, string) (Card, error)
	GetWatchListedCards(string) ([]Card, error)
	GetAllCardsInCollection(string) ([]Card, error)
	AddCardToCollection(Card, string) error
	UpdateCardInCollection(Card, string) error
}
