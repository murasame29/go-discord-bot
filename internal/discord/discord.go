package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/casino-bot/internal/discord/handler"
	"github.com/murasame29/casino-bot/internal/repository/mock"
)

type discord struct {
	token string
}

func New(token string) *discord {
	return &discord{
		token: token,
	}
}

func (d *discord) Start() error {
	dg, err := discordgo.New("Bot " + d.token)
	if err != nil {
		return err
	}

	userRepo := mock.NewUserRepo()
	gameRepo := mock.NewGameRepo()

	// handler を呼び出す
	h := handler.New(dg, userRepo, gameRepo)

	// ユーザーの登録
	h.SetupUsers()
	// BlackJackの登録
	h.Blackjack()

	// Open a websocket connection to Discord and begin listening.
	if err = dg.Open(); err != nil {
		return err
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return dg.Close()
}
