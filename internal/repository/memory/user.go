package memory

import (
	"context"

	"github.com/murasame29/casino-bot/internal/models"
	"github.com/murasame29/casino-bot/internal/repository"
)

type userRepo struct {
	users map[string]models.User
}

// NewUserRepo returns a new instance of userRepo.
func NewUserRepo() repository.UserRepo {
	return &userRepo{
		users: make(map[string]models.User),
	}
}

func (r *userRepo) Create(ctx context.Context, user models.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *userRepo) Get(ctx context.Context, id string) (*models.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, models.ErrUserNotFound
	}
	return &user, nil
}

func (r *userRepo) AddBalance(ctx context.Context, id string, amount int64) error {
	user, ok := r.users[id]
	if !ok {
		return models.ErrUserNotFound
	}
	user.Balance += amount
	r.users[id] = user
	return nil
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	delete(r.users, id)
	return nil
}
