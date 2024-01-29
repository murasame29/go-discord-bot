package mock

import (
	"github.com/murasame29/casino-bot/internal/models"
	"github.com/murasame29/casino-bot/internal/repository"
)

type gameRepo struct {
	games map[string]models.BlackJack
}

func NewGameRepo() repository.GameRepo {
	return &gameRepo{
		games: make(map[string]models.BlackJack),
	}
}

func (r *gameRepo) Create(game models.BlackJack) error {
	r.games[game.ID] = game
	return nil
}

func (r *gameRepo) Get(id string) (*models.BlackJack, error) {
	game, ok := r.games[id]
	if !ok {
		return nil, models.ErrGameNotFound
	}
	return &game, nil
}

func (r *gameRepo) Update(game models.BlackJack) error {
	r.games[game.ID] = game
	return nil
}

func (r *gameRepo) Delete(id string) error {
	delete(r.games, id)
	return nil
}
