package cardservice

// Card is a wrapper struct with Base card and Insert
type Card struct {
	Base
	Insert
}

// Base has fields belonging to a base card
type Base struct {
	Year   uint   `json:"year"`
	Set    string `json:"set"`
	Make   string `json:"make"`
	Number uint   `json:"number"`
	Player string `json:"player"`
}

// Insert struct holds extra values for the card
type Insert struct {
	Type        string `json:"type"`
	Memorabilia string `json:"memorabilia"`
	NumberedTo  string `json:"numberedTo"`
}
