package card

import (
	"reflect"
	"testing"
)

func TestNewDeck(t *testing.T) {
	cards := NewDeck()

	if len(cards) != 54 {
		t.Errorf("Expected 52 cards, got %d", len(cards))
	}

	for _, card := range cards {
		if card.IsJoker() {
			continue
		}

		// Check if card is in the correct suite
		if card.Rank() < Ace || card.Rank() > King {
			t.Errorf("Expected rank between Ace and King, got %d", card.Rank())
		}

		// Check if card is in the correct rank
		if card.Suit() != Spades && card.Suit() != Hearts && card.Suit() != Diamonds && card.Suit() != Clubs {
			t.Errorf("Expected suite to be Spades, Hearts, Diamonds or Clubs, got %s", card.Suit())
		}

		// Check if card is in the correct rank
		if rank[card.Rank()] != card.RankString() {
			t.Errorf("Expected rank string to be %s, got %s", rank[card.Rank()], card.RankString())
		}
	}

	// Check if there are 2 jokers
	jokerCount := 0
	for _, card := range cards {
		if card.IsJoker() {
			jokerCount++
		}
	}

	if jokerCount != 2 {
		t.Errorf("Expected 2 jokers, got %d", jokerCount)
	}

	// shuffle
	shuffled := Shuffle(cards)

	// Check if the cards are shuffled
	if reflect.DeepEqual(cards, shuffled) {
		t.Errorf("Expected cards to be shuffled, got %v", shuffled)
	}
}
