package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/casino-bot/internal/game/bj"
	"github.com/murasame29/casino-bot/internal/repository"
)

type handler struct {
	dg       *discordgo.Session
	userRepo repository.UserRepo
	gameRepo repository.BjRepo
	bj       bj.Game
}

// New returns a new instance of handler.
func New(
	dg *discordgo.Session,
	userRepo repository.UserRepo,
	gameRepo repository.BjRepo) *handler {
	return &handler{
		dg:       dg,
		userRepo: userRepo,
		gameRepo: gameRepo,
		bj:       bj.NewGame(gameRepo, userRepo),
	}
}
