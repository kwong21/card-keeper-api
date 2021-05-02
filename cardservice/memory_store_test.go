package cardservice

import (
	"fmt"
	"reflect"
	"testing"
)

var memoryStoreTest Repository

func TestMemoryStore(t *testing.T) {
	memoryStore, err := InMemoryStore()

	if err != nil {
		t.FailNow()
	}

	base := Base{
		Year:   "2020",
		Make:   "Upper Deck",
		Set:    "Series One",
		Player: "Brock Boeser",
	}

	fixtureCard := Card{
		Base:          base,
		IsOnWatchList: false,
	}

	fixtureCard.setCardID()

	t.Run(fmt.Sprintln("Test add card to collection"), func(t *testing.T) {
		err := memoryStore.AddCardToCollection(fixtureCard, "hockey")
		cards, err := memoryStore.GetAllCardsInCollection("hockey")

		if err != nil {
			t.Error("Failed adding card to collection")
		}

		if len(cards) != 1 {
			t.Error(fmt.Printf("Expected to find 1 card in collection, got %d", len(cards)))
		}
	})

	t.Run(fmt.Sprintf("Test get card from collection"), func(t *testing.T) {
		cardID := fixtureCard.CardID

		card, err := memoryStore.GetCardInCollection(cardID, "hockey")

		if !reflect.DeepEqual(card, fixtureCard) || err != nil {
			t.Error("Expected to retrieve fixture card from in-memory store")
		}
	})

	t.Run(fmt.Sprintf("Get Cards in watchlist"), func(t *testing.T) {
		base := Base{
			Year:   "2020",
			Make:   "Upper Deck",
			Set:    "Series One",
			Player: "Mike Sillinger",
		}

		watchListedCard := Card{
			Base:          base,
			IsOnWatchList: true,
		}

		watchListedCard.setCardID()

		err := memoryStore.AddCardToCollection(watchListedCard, "hockey")

		if err != nil {
			t.Error(fmt.Sprintf("Error adding watched listed card to collection: %s", err))
			t.FailNow()
		}

		watchListedCards, err := memoryStore.GetWatchListedCards("hockey")

		if err != nil {
			t.Error("Error getting watched listed cards")
		}

		if len(watchListedCards) != 1 {
			t.Fail()
		}

		c := watchListedCards[0]

		if c.CardID != watchListedCard.CardID {
			t.Error("Did not get watch listed card back from memoryStore")
		}
	})

	t.Run(fmt.Sprintln("Test should fail when adding duplicate card"), func(t *testing.T) {
		err := memoryStore.AddCardToCollection(fixtureCard, "hockey")

		switch et := err.(type) {
		default:
			t.Error(fmt.Sprintf("Expected duplicate error got %s", et.Error()))
		case nil:
			t.Error(fmt.Sprintf("Expected an error, but got nil"))
		case *DuplicateError:
		}
	})

	t.Run(fmt.Sprintln(), func(t *testing.T) {
		// err := memoryStore.UpdateCardInCollection(fixtureCard, "hockey")
	})
}
