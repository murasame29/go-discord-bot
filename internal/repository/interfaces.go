package repository

import (
	"context"

	"github.com/murasame29/casino-bot/internal/models"
)

type UserRepo interface {
	Create(ctx context.Context, user models.User) error
	Get(ctx context.Context, id string) (*models.User, error)
	AddBalance(ctx context.Context, id string, amount int64) error
	Delete(ctx context.Context, id string) error
}

type BjRepo interface {
	Create(ctx context.Context, game models.BlackJack) error
	Get(ctx context.Context, id string) (*models.BlackJack, error)
	Update(ctx context.Context, game models.BlackJack) error
	Delete(ctx context.Context, id string) error
}
