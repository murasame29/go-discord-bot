package main

import (
	"os"

	"github.com/murasame29/casino-bot/internal/discord"
)

var Token string

func init() {
	Token = os.Getenv("DISCORD_TOKEN")
}

func main() {
	d := discord.New(Token)
	if err := d.Start(); err != nil {
		panic(err)
	}
}
