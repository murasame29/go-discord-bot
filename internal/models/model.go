package models

import (
	"github.com/murasame29/casino-bot/internal/deck"
	"github.com/murasame29/casino-bot/internal/game/bj/hand"
)

type User struct {
	ID          string
	DisplayName string
	Balance     int64
}

type BlackJack struct {
	ID         string
	UserID     string
	Deck       deck.Deck
	DealerHand hand.Hand
	UserHand   []hand.Hand
	BetAmount  int64
	Insurance  int64
}
