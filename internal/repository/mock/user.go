package mock

import (
	"github.com/murasame29/casino-bot/internal/models"
	"github.com/murasame29/casino-bot/internal/repository"
)

type userRepo struct {
	users map[string]models.User
}

func NewUserRepo() repository.UserRepo {
	return &userRepo{
		users: make(map[string]models.User),
	}
}

func (r *userRepo) Create(user models.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *userRepo) Get(id string) (*models.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, models.ErrUserNotFound
	}
	return &user, nil
}

func (r *userRepo) AddBalance(id string, amount int64) error {
	user, ok := r.users[id]
	if !ok {
		return models.ErrUserNotFound
	}
	user.Balance += amount
	r.users[id] = user
	return nil
}

func (r *userRepo) Delete(id string) error {
	delete(r.users, id)
	return nil
}
