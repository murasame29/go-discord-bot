package handler

import "github.com/bwmarrin/discordgo"

// HelpMessage is the help message
const HelpMessage = `
**!register** - Botに対してユーザの登録を行います
**!balance** - ユーザの残高を確認します
**!bj [金額]** - ブラックジャックを開始します
**!bj-stand** - ブラックジャックでスタンドします
**!bj-hit** - ブラックジャックでヒットします
**!bj-dd** - ブラックジャックでダブルダウンします
**!bj-split** - ブラックジャックでスプリットします
**!bj-sur** - ブラックジャックでサレンダーします
**!bj-ins** - ブラックジャックでインシュランスします
**!az-vision [img]** - Azureの画像認識を行います
**!az-tr [from:to]** - Azureの翻訳を行います
`

func (h *handler) Help() {
	h.dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "!help" {
			s.ChannelMessageSend(m.ChannelID, HelpMessage)
		}
	})
}
