package cardservice

// Service for interacting with data store
type Service struct {
	Repository
}

// Repository writes and reads data from declared store
type Repository interface {
	GetAll() (*[]Card, error)
	AddCard(Card) error
}
