package handler

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/casino-bot/internal/models"
)

func (h *handler) SetupUsers() {
	h.dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "!register" {
			guild, err := s.Guild(m.GuildID)
			if err != nil {
				log.Println(err)
				return
			}

			for _, member := range guild.Members {
				// Add member to database
				err := h.userRepo.Create(models.User{
					ID:          member.User.ID,
					DisplayName: member.User.Username,
					Balance:     1000,
				})

				log.Println(member.User.ID)

				if err != nil {
					log.Println(err)
					return
				}
			}

			s.ChannelMessageSend(m.ChannelID, "all users added to database")
		}
	})
}
