package deck

import "testing"

func TestDeck(t *testing.T) {
	deck := New(1)

	if deck.IsEmpty() {
		t.Errorf("Expected deck to not be empty")
	}

	deck.Shuffle()

	if deck.IsEmpty() {
		t.Errorf("Expected deck to not be empty")
	}

	for !deck.IsEmpty() {
		deck.Draw()
	}

	if !deck.IsEmpty() {
		t.Errorf("Expected deck to be empty")
	}
}

func TestDeckWithoutJokers(t *testing.T) {
	deck := New(1, IgnoreJokers())

	if deck.IsEmpty() {
		t.Errorf("Expected deck to not be empty")
	}

	for !deck.IsEmpty() {
		if deck.Draw().IsJoker() {
			t.Errorf("Expected deck to not have jokers")
		}
	}
}

func TestDeckSize(t *testing.T) {
	deck := New(2)

	cnt := 0
	for !deck.IsEmpty() {
		deck.Draw()
		cnt++
	}

	if cnt != 108 {
		t.Errorf("Expected deck to have 108 cards, got %d", cnt)
	}
}
