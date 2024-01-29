package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/casino-bot/internal/game/bj/hand"
	"github.com/murasame29/casino-bot/internal/models"
)

const (
	CommandStart = "!bj"
	CommandHit   = "!hit"
	CommandStand = "!stand"
	CommandSplit = "!split"
	CommandDD    = "!dd"
	CommandIns   = "!ins"
	CommandSur   = "!sur"
)

func (h *handler) Blackjack() {
	h.dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		switch strings.Split(m.Content, " ")[0] {
		case CommandStart:
			h.startGame(s, m)
		case CommandHit:
			h.hit(s, m)
		case CommandStand:
			h.stand(s, m)
		case CommandSplit:
			h.split(s, m)
		case CommandDD:
			h.doubleDown(s, m)
		case CommandIns:
			h.insurance(s, m)
		case CommandSur:
			h.surrender(s, m)
		default:
			return
		}
	})
}

func (h *handler) startGame(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if _, err := h.userRepo.Get(m.Author.ID); err != nil {
		s.ChannelMessageSend(m.ChannelID, "初めて見る顔ですね。勝手に登録しておきます。所持金は1000です。")
		if err := h.userRepo.Create(models.User{
			ID:          m.Author.ID,
			DisplayName: m.Author.Username,
			Balance:     1000,
		}); err != nil {
			s.ChannelMessageSend(m.ChannelID, "登録に失敗しました。")
			return
		}

		s.ChannelMessageSend(m.ChannelID, "登録が完了しました。")
	}

	s.ChannelMessageSend(m.ChannelID, "BlackJackを開始します。")

	args := strings.Split(m.Content, " ")

	if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "コマンドが誤っています \n 100を賭ける場合は !bh 100 と入力してください。")
		return
	}

	betAmount, err := strconv.Atoi(args[1])
	if err != nil || betAmount < 1 {
		s.ChannelMessageSend(m.ChannelID, "賭け金が不正です。 1以上の整数を入力してください。")
		return
	}

	out, err := h.bj.Start(m.Author.ID, int64(betAmount))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var (
		hands         []string
		userHandValue int
	)
	// 開始早々BJの場合
	if out.IsEnd {
		for _, hand := range out.GameData.UserHand {
			hands = append(hands, strings.Join(hand.Strings(), ", "))
			userHandValue += hand.Score()
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(strings.Join(StartGameMessage, "\n"), strings.Join(hands, ", "), userHandValue, out.UserData.Balance))
		return
	}

	hands = append(hands, fmt.Sprintf("%s", strings.Join(out.GameData.UserHand[0].Strings(), ", ")))
	userHandValue += out.GameData.UserHand[0].Score()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(strings.Join(StartGameMessage, "\n"), out.GameData.BetAmount, out.UserData.Balance, out.GameData.DealerHand.Strings()[0], out.GameData.DealerHand.RawCards()[0].BJscore(), strings.Join(hands, ", "), userHandValue))
}

func (h *handler) hit(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	out, err := h.bj.Hit(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var (
		hands         [][]string
		userHandValue []int
	)

	// ゲームが終了した場合 (burst)
	if out.IsEnd {
		for _, hand := range out.GameData.UserHand {
			hands = append(hands, hand.Strings())
			userHandValue = append(userHandValue, hand.Score())
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(BustMessage, strings.Join(hands[0], ", "), userHandValue[0], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
		return
	}

	// Splitされている場合
	if len(out.GameData.UserHand) > 1 {
		for _, hand := range out.GameData.UserHand {
			hands = append(hands, hand.Strings())
			userHandValue = append(userHandValue, hand.Score())
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(SplitedHitMessage, strings.Join(hands[0], ", "), strings.Join(hands[1], ", "), userHandValue[0], userHandValue[1]))
		return
	}

	hands = append(hands, out.GameData.UserHand[0].Strings())
	userHandValue = append(userHandValue, out.GameData.UserHand[0].Score())

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(HitMessage, strings.Join(hands[0], ", "), userHandValue[0]))
}

func (h *handler) stand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	var (
		handID int = 1
		err    error
	)
	args := strings.Split(m.Content, " ")
	if len(args) == 2 {
		handID, err = strconv.Atoi(args[1])
	}

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "コマンドが誤っています \n 2つ目の手札をスタンドする場合は !stand 2 と入力してください。")
		return
	}

	out, err := h.bj.Stand(m.Author.ID, handID-1)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var (
		hands         [][]string
		userHandValue []int
	)

	// Splitされている場合
	if len(out.GameData.UserHand) > 1 {
		// ゲームが終了している場合
		if out.IsEnd {
			for _, uhand := range out.GameData.UserHand {
				hands = append(hands, uhand.Strings())
				userHandValue = append(userHandValue, uhand.Score())

				// 勝敗判定
				switch uhand.Status() {
				case hand.StatusWin:
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(SplitedWinMessage, strings.Join(hands[0], ", "), strings.Join(hands[1], ", "), userHandValue[0], userHandValue[1], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
				case hand.StatusLose:
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(SplitedLoseMessage, strings.Join(hands[0], ", "), strings.Join(hands[1], ", "), userHandValue[0], userHandValue[1], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
				case hand.StatusDraw:
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(SplitedDrawMessage, strings.Join(hands[0], ", "), strings.Join(hands[1], ", "), userHandValue[0], userHandValue[1], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
				}
			}
			return
		}

		// ゲームが終了していない場合
		for _, hand := range out.GameData.UserHand {
			hands = append(hands, hand.Strings())
			userHandValue = append(userHandValue, hand.Score())
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(SplitedStandMessage, strings.Join(hands[0], ", "), strings.Join(hands[1], ", "), userHandValue[0], userHandValue[1]))
		return
	}

	// ゲームが終了している場合
	if out.IsEnd {
		// インシュランスに勝った場合
		if out.IsInsuranceWin && out.GameData.Insurance != 0 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(InsuranceWinMessage, strings.Join(out.GameData.UserHand[0].Strings(), ", "), out.GameData.UserHand[0].Score(), out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
		}

		for _, hand := range out.GameData.UserHand {
			hands = append(hands, hand.Strings())
			userHandValue = append(userHandValue, hand.Score())
		}

		// 勝敗判定
		switch out.GameData.UserHand[0].Status() {
		case hand.StatusWin:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(WinMessage, strings.Join(hands[0], ", "), userHandValue[0], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
		case hand.StatusLose:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(LoseMessage, strings.Join(hands[0], ", "), userHandValue[0], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
		case hand.StatusDraw:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(DrawMessage, strings.Join(hands[0], ", "), userHandValue[0], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
		}
		return
	}
	hands = append(hands, out.GameData.UserHand[0].Strings())
	userHandValue = append(userHandValue, out.GameData.UserHand[0].Score())

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(StandMessage, strings.Join(hands[0], ", "), userHandValue[0]))
}

func (h *handler) split(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	out, err := h.bj.Split(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var (
		hands         [][]string
		userHandValue []int
	)

	for _, hand := range out.GameData.UserHand {
		hands = append(hands, hand.Strings())
		userHandValue = append(userHandValue, hand.Score())
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(SplitMessage, strings.Join(hands[0], ", "), strings.Join(hands[1], ", "), userHandValue[0], userHandValue[1]))
}

func (h *handler) doubleDown(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	out, err := h.bj.DoubleDown(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var (
		hands         []string
		userHandValue int
	)

	for _, hand := range out.GameData.UserHand {
		hands = append(hands, strings.Join(hand.Strings(), ", "))
		userHandValue += hand.Score()
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(DoubleDownMessage, strings.Join(hands, ", "), userHandValue))
}

func (h *handler) insurance(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	args := strings.Split(m.Content, " ")

	if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "コマンドが誤っています \n 100を賭ける場合は !bh 100 と入力してください。")
		return
	}

	betAmount, err := strconv.Atoi(args[1])
	if err != nil || betAmount < 1 {
		s.ChannelMessageSend(m.ChannelID, "賭け金が不正です。 1以上の整数を入力してください。")
		return
	}

	out, err := h.bj.Insurance(m.Author.ID, int64(betAmount))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var (
		hands         []string
		userHandValue int
	)

	for _, hand := range out.GameData.UserHand {
		hands = append(hands, strings.Join(hand.Strings(), ", "))
		userHandValue += hand.Score()
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(InsuranceMessage, strings.Join(hands, ", "), userHandValue, betAmount, out.UserData.Balance))
}

func (h *handler) surrender(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	out, err := h.bj.Surrender(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var (
		hands         []string
		userHandValue int
	)

	for _, hand := range out.GameData.UserHand {
		hands = append(hands, strings.Join(hand.Strings(), ", "))
		userHandValue += hand.Score()
	}

	s.ChannelMessageSend(m.ChannelID, SurrenderMessage)
}
