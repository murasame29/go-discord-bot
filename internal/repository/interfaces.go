package repository

import "github.com/murasame29/casino-bot/internal/models"

type UserRepo interface {
	Create(user models.User) error
	Get(id string) (*models.User, error)
	AddBalance(id string, amount int64) error
	Delete(id string) error
}

type GameRepo interface {
	Create(game models.BlackJack) error
	Get(id string) (*models.BlackJack, error)
	Update(game models.BlackJack) error
	Delete(id string) error
}
