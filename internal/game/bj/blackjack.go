package bj

import (
	"github.com/murasame29/casino-bot/internal/deck"
	"github.com/murasame29/casino-bot/internal/game/bj/hand"
	"github.com/murasame29/casino-bot/internal/models"
	"github.com/murasame29/casino-bot/internal/repository"
)

type game struct {
	gameRepo repository.GameRepo
	userRepo repository.UserRepo
}

type Game interface {
	Start(userID string, betAmount int64) (*OutGame, error)
	Hit(userID string) (*OutGame, error)
	Stand(userID string, handID int) (*OutGame, error)
	DoubleDown(userID string) (*OutGame, error)
	Split(userID string) (*OutGame, error)
	Surrender(userID string) (*OutGame, error)
	Insurance(userID string, insurance int64) (*OutGame, error)
}

func NewGame(gameRepo repository.GameRepo, userRepo repository.UserRepo) Game {
	return &game{
		gameRepo: gameRepo,
		userRepo: userRepo,
	}
}

func (g *game) Start(userID string, betAmount int64) (*OutGame, error) {
	user, err := g.userRepo.Get(userID)
	if err != nil {
		return nil, err
	}

	// 所持金より高い額を賭けられない
	if user.Balance < betAmount {
		return nil, models.ErrInsufficientBalance
	}

	// ゲームが進行中の場合はエラー
	if _, err = g.gameRepo.Get(userID); err != models.ErrGameNotFound {
		return nil, models.ErrGameDuplicate
	}

	// ゲームを作成
	game := models.BlackJack{
		ID:         userID,
		UserID:     userID,
		Deck:       deck.New(1, deck.IgnoreJokers()),
		DealerHand: hand.NewHand(),
		UserHand:   []hand.Hand{hand.NewHand()},
		BetAmount:  betAmount,
	}

	// ディーラーとプレイヤーにカードを配る
	for i := 0; i < 2; i++ {
		game.DealerHand.Add(game.Deck.Draw())
		game.UserHand[0].Add(game.Deck.Draw())
	}

	// PlayerがBlackJackの場合はゲームを終了する
	if game.UserHand[0].IsBlackJack() {
		// ゲームを削除
		if err := g.gameRepo.Delete(userID); err != nil {
			return nil, err
		}

		// ユーザーの所持金を更新
		if err := g.userRepo.AddBalance(userID, betAmount*2); err != nil {
			return nil, err
		}

		return &OutGame{
			GameData: &game,
			UserData: user,
			IsEnd:    true,
		}, nil
	}

	// ゲームを保存
	if err := g.gameRepo.Create(game); err != nil {
		return nil, err
	}

	// ユーザーの所持金を更新
	if err := g.userRepo.AddBalance(userID, -betAmount); err != nil {
		return nil, err
	}

	return &OutGame{
		GameData: &game,
		UserData: user,
		IsEnd:    false,
	}, nil
}

func (g *game) Hit(userID string) (*OutGame, error) {
	user, err := g.userRepo.Get(userID)
	if err != nil {
		return nil, err
	}

	game, err := g.gameRepo.Get(userID)
	if err != nil {
		return nil, models.ErrGameNotAvailable
	}

	// カードを引く
	for i := 0; i < len(game.UserHand); i++ {
		if game.UserHand[i].IsStand() {
			continue
		}

		game.UserHand[i].Add(game.Deck.Draw())
	}

	// Splitしてない状態でBustしたらゲームを終了する
	if len(game.UserHand) == 1 && game.UserHand[0].IsBust() {
		// ゲームを削除
		if err := g.gameRepo.Delete(userID); err != nil {
			return nil, err
		}

		return &OutGame{
			GameData: game,
			UserData: user,
			IsEnd:    true,
		}, nil
	}

	// Splitしている場合は全ての手札がBustしていないか確認する
	if len(game.UserHand) == 2 {
		var isAllBust bool
		for _, hand := range game.UserHand {
			if !hand.IsBust() {
				isAllBust = false
				break
			}
			isAllBust = true
		}

		if isAllBust {
			// ゲームを削除
			if err := g.gameRepo.Delete(userID); err != nil {
				return nil, err
			}

			return &OutGame{
				GameData: game,
				UserData: user,
				IsEnd:    true,
			}, nil
		}
	}

	// ゲームを保存
	if err := g.gameRepo.Update(*game); err != nil {
		return nil, err
	}

	return &OutGame{
		GameData: game,
		UserData: user,
		IsEnd:    false,
	}, nil
}

func (g *game) Stand(userID string, handID int) (*OutGame, error) {
	user, err := g.userRepo.Get(userID)
	if err != nil {
		return nil, err
	}

	game, err := g.gameRepo.Get(userID)
	if err != nil {
		return nil, models.ErrGameNotAvailable
	}

	// 0 or 1 若しくはスプリットしていない場合はエラー
	if handID != 0 && handID != 1 || len(game.UserHand) != 1 && handID == 1 {
		return nil, models.ErrBadCommand
	}
	// 手札をStandにする
	game.UserHand[handID].Stand()

	// 全ての手札がStandになっていない場合はゲームを続行する
	if !checkAllHandStand(*game) {
		// ゲームを保存
		if err := g.gameRepo.Update(*game); err != nil {
			return nil, err
		}
		return &OutGame{
			GameData: game,
			UserData: user,
		}, nil
	}

	// DealerのHandがBlackJackの場合はゲームを終了する
	if game.DealerHand.IsBlackJack() {
		// 保険があれば支払う
		if game.Insurance != 0 {
			// ユーザーの所持金を更新
			if err := g.userRepo.AddBalance(userID, game.Insurance*3); err != nil {
				return nil, err
			}
		}
		// ゲームを削除
		if err := g.gameRepo.Delete(userID); err != nil {
			return nil, err
		}
		return &OutGame{
			GameData:       game,
			UserData:       user,
			IsEnd:          true,
			IsInsuranceWin: true,
		}, nil
	}

	// DealerのHandが17以上になるまでカードを引く
	for game.DealerHand.Score() < 17 {
		game.DealerHand.Add(game.Deck.Draw())
	}

	// DealerのHandが21を超えたらBust
	if game.DealerHand.IsBust() {
		// ゲームを保存
		if err := g.gameRepo.Update(*game); err != nil {
			return nil, err
		}
		// スプリットしてたらベット額を倍にする
		if len(game.UserHand) == 2 {
			game.BetAmount *= 2
		}
		// ユーザーの所持金を更新
		if err := g.userRepo.AddBalance(userID, game.BetAmount*2); err != nil {
			return nil, err
		}
		return &OutGame{
			GameData: game,
			UserData: user,
			IsEnd:    true,
		}, nil
	}

	// DealerのHandが21以下の場合、ユーザーのHandと比較する
	for _, uhand := range game.UserHand {
		// 既にBustしている場合はスキップ
		if uhand.Status() > hand.StatusStand {
			continue
		}
		// DealerのHandがユーザーのHandより大きい場合は負け
		if game.DealerHand.Score() > uhand.Score() {
			uhand.UpdateStatus(hand.StatusLose)
		} else if game.DealerHand.Score() < uhand.Score() {
			uhand.UpdateStatus(hand.StatusWin)
			// 所持金更新
			if err := g.userRepo.AddBalance(userID, game.BetAmount*2); err != nil {
				return nil, err
			}
		}
	}

	// ゲームを削除
	if err := g.gameRepo.Delete(userID); err != nil {
		return nil, err
	}

	return &OutGame{
		GameData: game,
		UserData: user,
		IsEnd:    true,
	}, nil
}

func (g *game) DoubleDown(userID string) (*OutGame, error) {
	user, err := g.userRepo.Get(userID)
	if err != nil {
		return nil, err
	}

	game, err := g.gameRepo.Get(userID)
	if err != nil {
		return nil, models.ErrGameNotAvailable
	}

	// 所持金より高い額を賭けられない
	if user.Balance < game.BetAmount {
		return nil, models.ErrInsufficientBalance
	}

	// Splitされている場合はDoubleDownできない
	if len(game.UserHand) != 1 {
		return nil, models.ErrBadCommand
	}

	// カードを引く
	game.UserHand[0].Add(game.Deck.Draw())

	// ゲームを保存
	if err := g.gameRepo.Update(*game); err != nil {
		return nil, err
	}

	// ユーザーの所持金を更新
	if err := g.userRepo.AddBalance(userID, -game.BetAmount); err != nil {
		return nil, err
	}

	return g.Stand(userID, 0)
}

func (g *game) Split(userID string) (*OutGame, error) {
	user, err := g.userRepo.Get(userID)
	if err != nil {
		return nil, err
	}

	game, err := g.gameRepo.Get(userID)
	if err != nil {
		return nil, models.ErrGameNotAvailable
	}

	// 所持金より高い額を賭けられない
	if user.Balance < game.BetAmount {
		return nil, models.ErrInsufficientBalance
	}

	// Splitできるのは2枚のカードが同じランク且つスプリットしていない場合且つドローされてない状態
	if !game.UserHand[0].IsPair() || len(game.UserHand) != 1 || game.UserHand[0].Len() != 2 {
		return nil, models.ErrBadCommand
	}

	// 手札を分割する
	game.UserHand = game.UserHand[0].SplitHand()

	// ゲームを保存
	if err := g.gameRepo.Update(*game); err != nil {
		return nil, err
	}

	// ユーザーの所持金を更新
	if err := g.userRepo.AddBalance(userID, -game.BetAmount); err != nil {
		return nil, err
	}

	return &OutGame{
		GameData: game,
		UserData: user,
	}, nil
}

func (g *game) Surrender(userID string) (*OutGame, error) {
	user, err := g.userRepo.Get(userID)
	if err != nil {
		return nil, err
	}

	game, err := g.gameRepo.Get(userID)
	if err != nil {
		return nil, models.ErrGameNotAvailable
	}

	// ゲームを削除
	if err := g.gameRepo.Delete(userID); err != nil {
		return nil, err
	}

	// ユーザーの所持金を更新
	if err := g.userRepo.AddBalance(userID, game.BetAmount/2); err != nil {
		return nil, err
	}

	return &OutGame{
		GameData: game,
		UserData: user,
	}, nil
}

func (g *game) Insurance(userID string, insurance int64) (*OutGame, error) {
	user, err := g.userRepo.Get(userID)
	if err != nil {
		return nil, err
	}

	game, err := g.gameRepo.Get(userID)
	if err != nil {
		return nil, models.ErrGameNotAvailable
	}

	// 所持金より高い額を賭けられない
	if user.Balance < insurance {
		return nil, models.ErrInsufficientBalance
	}

	// Insuranceは最初のカードがAの場合のみ
	if game.DealerHand.RawCards()[0].Rank() != 1 {
		return nil, models.ErrBadCommand
	}

	// 元の金額の半額を超えたらエラー
	if insurance > game.BetAmount/2 {
		return nil, models.ErrBadCommand
	}

	game.Insurance = insurance

	// ゲームを保存
	if err := g.gameRepo.Update(*game); err != nil {
		return nil, err
	}

	// ユーザーの所持金を更新
	if err := g.userRepo.AddBalance(userID, -insurance); err != nil {
		return nil, err
	}

	user.Balance -= insurance

	return &OutGame{
		GameData: game,
		UserData: user,
		IsEnd:    false,
	}, nil
}

func checkAllHandStand(game models.BlackJack) bool {
	for _, hand := range game.UserHand {
		if !hand.IsStand() {
			return false
		}
	}
	return true
}
