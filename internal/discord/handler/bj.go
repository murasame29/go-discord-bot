package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/casino-bot/internal/game/bj/hand"
	"github.com/murasame29/casino-bot/internal/models"
)

const (
	// CommandStart is a message that is sent when the game starts.
	CommandStart = "!bj"
	CommandHit   = "!bj-hit"
	CommandStand = "!bj-stand"
	CommandSplit = "!bj-split"
	CommandDD    = "!bj-dd"
	CommandIns   = "!bj-ins"
	CommandSur   = "!bj-sur"
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

	ctx := context.Background()

	th, err := s.MessageThreadStart(m.ChannelID, m.ID, fmt.Sprintf("BlackJack-%s", m.Author.ID), 0)
	if err != nil {
		log.Println(err)
	}
	channelID := th.ID

	if _, err := h.userRepo.Get(ctx, m.Author.ID); err != nil {
		s.ChannelMessageSend(channelID, "初めて見る顔ですね。勝手に登録しておきます。所持金は1000です。")
		if err := h.userRepo.Create(ctx, models.User{
			ID:          m.Author.ID,
			DisplayName: m.Author.Username,
			Balance:     1000,
		}); err != nil {
			s.ChannelMessageSend(channelID, "登録に失敗しました。")
			return
		}

		s.ChannelMessageSend(channelID, "登録が完了しました。")
	}

	s.ChannelMessageSend(channelID, "BlackJackを開始します。")

	args := strings.Split(m.Content, " ")

	if len(args) != 2 {
		s.ChannelMessageSend(channelID, "コマンドが誤っています \n 100を賭ける場合は !bh 100 と入力してください。")
		return
	}

	betAmount, err := strconv.Atoi(args[1])
	if err != nil || betAmount < 1 {
		s.ChannelMessageSend(channelID, "賭け金が不正です。 1以上の整数を入力してください。")
		return
	}

	out, err := h.bj.Start(ctx, m.Author.ID, int64(betAmount))
	if err != nil {
		s.ChannelMessageSend(channelID, err.Error())
		return
	}

	var (
		hands         []string
		userHandValue int
	)

	hands = append(hands, fmt.Sprintf("%s", strings.Join(out.GameData.UserHand[0].Strings(), ", ")))
	userHandValue += out.GameData.UserHand[0].Score()

	s.ChannelMessageSend(channelID, fmt.Sprintf(strings.Join(StartGameMessage, "\n"), out.GameData.BetAmount, out.UserData.Balance, out.GameData.DealerHand.Strings()[0], out.GameData.DealerHand.RawCards()[0].BJscore(), strings.Join(hands, ", "), userHandValue))

	// 開始早々BJの場合
	if out.IsEnd {
		s.ChannelMessageSend(channelID, fmt.Sprintf(strings.Join(BJMessage, "\n"), out.UserData.Balance))
		h.deleteChannel(s, channelID)
	} else {
		s.ChannelMessageSend(channelID, NextStepMessage)
	}
}

func (h *handler) hit(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	ctx := context.Background()

	out, err := h.bj.Hit(ctx, m.Author.ID)
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
		h.deleteChannel(s, m.ChannelID)
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

	ctx := context.Background()
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

	out, err := h.bj.Stand(ctx, m.Author.ID, handID-1)
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
			for i, uhand := range out.GameData.UserHand {
				hands = append(hands, uhand.Strings())
				userHandValue = append(userHandValue, uhand.Score())

				// 勝敗判定
				switch uhand.Status() {
				case hand.StatusWin:
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(SplitedWinMessage, i+1, strings.Join(hands[0], ", "), strings.Join(hands[1], ", "), userHandValue[0], userHandValue[1], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score()))
				case hand.StatusLose:
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(SplitedLoseMessage, i+1, strings.Join(hands[0], ", "), strings.Join(hands[1], ", "), userHandValue[0], userHandValue[1], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score()))
				case hand.StatusDraw:
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(SplitedDrawMessage, i+1, strings.Join(hands[0], ", "), strings.Join(hands[1], ", "), userHandValue[0], userHandValue[1], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score()))
				}
			}
			h.deleteChannel(s, m.ChannelID)
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
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(strings.Join(WinMessage, "\n"), strings.Join(hands[0], ", "), userHandValue[0], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
		case hand.StatusLose:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(strings.Join(LoseMessage, "\n"), strings.Join(hands[0], ", "), userHandValue[0], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
		case hand.StatusDraw:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(strings.Join(DrawMessage, "\n"), strings.Join(hands[0], ", "), userHandValue[0], out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
		}
		h.deleteChannel(s, m.ChannelID)
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
	ctx := context.Background()

	out, err := h.bj.Split(ctx, m.Author.ID)
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

	ctx := context.Background()
	out, err := h.bj.DoubleDown(ctx, m.Author.ID)
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

	switch out.GameData.UserHand[0].Status() {
	case hand.StatusWin:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(strings.Join(DoubleDownMessage, "\n"), strings.Join(hands, ", "), userHandValue, out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), WinResultMessage, out.UserData.Balance))
	case hand.StatusLose:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(strings.Join(DoubleDownMessage, "\n"), strings.Join(hands, ", "), userHandValue, out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), LoseResultMessage, out.UserData.Balance))
	case hand.StatusDraw:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(strings.Join(DoubleDownMessage, "\n"), strings.Join(hands, ", "), userHandValue, out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), DrawResultMessage, out.UserData.Balance))
	default:
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(strings.Join(DoubleDownMessage, "\n"), strings.Join(hands, ", "), userHandValue, out.GameData.DealerHand.Strings(), out.GameData.DealerHand.Score(), out.UserData.Balance))
	}
}

func (h *handler) insurance(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	ctx := context.Background()
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

	out, err := h.bj.Insurance(ctx, m.Author.ID, int64(betAmount))
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
	ctx := context.Background()
	out, err := h.bj.Surrender(ctx, m.Author.ID)
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
	h.deleteChannel(s, m.ChannelID)
}

func (h *handler) deleteChannel(s *discordgo.Session, id string) {
	time.Sleep(1 * time.Minute)
	s.ChannelDelete(id)
}
