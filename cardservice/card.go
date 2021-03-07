package cardservice

import (
	"hash/fnv"
	"strings"
)

// Card is a wrapper struct with Base card and Insert
type Card struct {
	CardID uint32 `bson:"cardID,omitempty"`
	Base
	Insert
	IsOnWatchList bool `json:"isWatchListed"`
}

// Base has fields belonging to a base card
type Base struct {
	Year   string `bson:"season,omitempty"  json:"season"`
	Set    string `bson:"set,omitempty" json:"set"`
	Make   string `bson:"manufacturer,omitempty" json:"manufacturer"`
	Number string `bson:"cardNumber,omitempty" json:"card_number"`
	Player string `bson:"player,omitempty" json:"player"`
}

// Insert struct holds extra values for the card
type Insert struct {
	Memorabilia string `json:"memorabilia"`
	NumberedTo  string `json:"numberedTo"`
}

func (c *Card) setCardID() {
	cardDetailString := c.Base.Year + c.Base.Set + c.Base.Make + c.Base.Number + c.Base.Player

	h := fnv.New32a()
	h.Write([]byte(strings.ToLower(cardDetailString)))

	c.CardID = h.Sum32()
}
