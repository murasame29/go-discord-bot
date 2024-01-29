package hand

import (
	"fmt"

	"github.com/murasame29/casino-bot/internal/deck/card"
)

type Status int

const (
	StatusNone Status = iota
	StatusPlaying
	StatusStand
	StatusWin
	StatusLose
	StatusDraw
)

type hand struct {
	cards  []card.Card
	status Status
}

type Hand interface {
	Add(card card.Card)
	Stand()

	IsBust() bool
	IsStand() bool
	IsBlackJack() bool
	IsPair() bool

	Score() int
	Len() int
	RawCards() []card.Card
	Strings() []string
	SplitHand() []Hand

	Status() Status
	UpdateStatus(status Status)
}

func NewHand() Hand {
	return &hand{
		cards:  make([]card.Card, 0),
		status: StatusPlaying,
	}
}

func (h *hand) Add(card card.Card) {
	h.cards = append(h.cards, card)
	if h.Score() > 21 {
		h.status = StatusLose
	}
}

func (h *hand) Stand() {
	h.status = StatusStand
}

func (h *hand) IsStand() bool {
	return h.status == StatusStand
}

func (h *hand) IsBust() bool {
	return h.status == StatusLose
}

func (h *hand) Len() int {
	return len(h.cards)
}

func (h *hand) RawCards() []card.Card {
	return h.cards
}

func (h *hand) Score() int {
	// Aがあれば一度後に回して計算
	var score int
	var AceCount int
	for _, card := range h.cards {
		if card.Rank() == 1 {
			AceCount++
			continue
		}
		score += card.BJscore()
	}

	for i := 0; i < AceCount; i++ {
		if score+11 > 21 {
			score++
		} else {
			score += 11
		}
	}

	return score
}

func (h *hand) UpdateStatus(status Status) {
	h.status = status
}

func (h *hand) Status() Status {
	return h.status
}

func (h *hand) Strings() []string {
	var hands []string
	for _, card := range h.cards {
		hands = append(hands, fmt.Sprintf("%s%d", card.Suit(), card.Rank()))
	}
	return hands
}

func (h *hand) SplitHand() []Hand {
	var hands []Hand
	for _, card := range h.cards {
		hands = append(hands, NewHand())
		hands[len(hands)-1].Add(card)
	}
	return hands
}

func (h *hand) IsBlackJack() bool {
	if len(h.cards) != 2 {
		return false
	}

	if h.cards[0].Rank() == 1 && h.cards[1].Rank() == 10 {
		return true
	}

	if h.cards[0].Rank() == 10 && h.cards[1].Rank() == 1 {
		return true
	}

	return false
}

func (h *hand) IsPair() bool {
	if len(h.cards) != 2 {
		return false
	}

	return h.cards[0].Rank() == h.cards[1].Rank()
}
