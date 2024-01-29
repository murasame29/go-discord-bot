package deck

import "github.com/murasame29/casino-bot/internal/deck/card"

type deck struct {
	cards []card.Card

	ignoreJokers bool
}

type Deck interface {
	Shuffle()
	Draw() card.Card
	IsEmpty() bool
}

type Option func(*deck)

func IgnoreJokers() Option {
	return func(d *deck) {
		d.ignoreJokers = true

		for i, card := range d.cards {
			if card.IsJoker() {
				d.cards = append(d.cards[:i], d.cards[i+1:]...)
			}
		}
	}
}

func New(deckSize int, opts ...Option) Deck {
	var cards []card.Card
	for i := 0; i < deckSize; i++ {
		cards = append(cards, card.NewDeck()...)
	}

	deck := &deck{
		cards: card.Shuffle(cards),
	}

	for _, opt := range opts {
		opt(deck)
	}

	return deck
}

func (d *deck) Shuffle() {
	d.cards = card.Shuffle(d.cards)
}

func (d *deck) Draw() card.Card {
	if d.IsEmpty() {
		return nil
	}

	card := d.cards[0]
	d.cards = d.cards[1:]

	return card
}

func (d *deck) IsEmpty() bool {
	return len(d.cards) == 0
}
