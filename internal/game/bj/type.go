package bj

import (
	"github.com/murasame29/casino-bot/internal/models"
)

type OutGame struct {
	GameData       *models.BlackJack
	UserData       *models.User
	IsEnd          bool
	IsInsuranceWin bool
}
