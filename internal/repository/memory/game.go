package memory

import (
	"context"

	"github.com/murasame29/casino-bot/internal/models"
	"github.com/murasame29/casino-bot/internal/repository"
)

type bjRepo struct {
	games map[string]models.BlackJack
}

// NewGameRepo returns a new instance of gameRepo.
func NewGameRepo() repository.BjRepo {
	return &bjRepo{
		games: make(map[string]models.BlackJack),
	}
}

func (r *bjRepo) Create(ctx context.Context, game models.BlackJack) error {
	r.games[game.ID] = game
	return nil
}

func (r *bjRepo) Get(ctx context.Context, id string) (*models.BlackJack, error) {
	game, ok := r.games[id]
	if !ok {
		return nil, models.ErrGameNotFound
	}
	return &game, nil
}

func (r *bjRepo) Update(ctx context.Context, game models.BlackJack) error {
	r.games[game.ID] = game
	return nil
}

func (r *bjRepo) Delete(ctx context.Context, id string) error {
	delete(r.games, id)
	return nil
}
