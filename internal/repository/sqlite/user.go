package sqlite

import (
	"context"
	"database/sql"

	"github.com/murasame29/casino-bot/internal/models"
	"github.com/murasame29/casino-bot/internal/repository"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) repository.UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user models.User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (id, display_name, balance) VALUES (?, ?, ?)", user.ID, user.DisplayName, user.Balance)
	return err
}

func (r *userRepo) Get(ctx context.Context, id string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, display_name, balance FROM users WHERE id = ?", id)
	var user models.User
	err := row.Scan(&user.ID, &user.DisplayName, &user.Balance)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) AddBalance(ctx context.Context, id string, amount int64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET balance = balance + ? WHERE id = ?", amount, id)
	return err
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	return err
}
