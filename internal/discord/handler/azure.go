package handler

import "github.com/bwmarrin/discordgo"

func (h *handler) azureAIs() {
	h.dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "!az" {
			s.ChannelMessageSend(m.ChannelID, "I am Azure, the AI of this server")
		}
	})
}
