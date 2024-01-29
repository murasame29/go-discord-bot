package handler

const (
	StartMessage = "%dP賭けてゲームを開始します。\nあなたの所持金は%dPです。\nディーラーのカードは [%s]です。 合計値は %d。\nあなたの手札は [%s] です。合計値は %d。\n次の操作を行ってください"
	BJMessage    = "おめでとうございます。ブラックジャックです。\nあなたの手札は [%s] です。合計値は %d。\nあなたの所持金は %dPです。"

	SplitedHitMessage   = "ヒットしました。\nあなたの手札は [%s] [%s] です。合計値は %dと%dです。\n次の操作を行ってください"
	SplitedStandMessage = "スタンドしました。\nあなたの手札は [%s] [%s] です。合計値は %dと%dです。\n次の操作を行ってください"

	SplitedWinMessage  = "%d枚目の手札で勝ちました。\nあなたの手札は [%s] [%s]です。合計値は %dと%dです。\nディーラーの手札は [%s] です。"
	SplitedLoseMessage = "%d枚目の手札で負けました。\nあなたの手札は [%s] [%s]です。合計値は %dと%dです。\nディーラーの手札は [%s] です。"
	SplitedDrawMessage = "%d枚目の手札で引き分けました。\nあなたの手札は [%s] [%s]です。合計値は %dと%dです。\nディーラーの手札は [%s] です。"

	HitMessage   = "ヒットしました。\nあなたの手札は [%s] です。合計値は %d。\n次の操作を行ってください"
	StandMessage = "スタンドしました。\nあなたの手札は [%s] です。合計値は %d。\n次の操作を行ってください"

	WinMessage  = "あなたの勝ちです。\nあなたの手札は [%s] です。合計値は %d。\nディーラーの手札は [%s] です。合計値は %d。\nあなたの所持金は %dPです。"
	LoseMessage = "あなたの負けです。\nあなたの手札は [%s] です。合計値は %d。\nディーラーの手札は [%s] です。合計値は %d。\nあなたの所持金は %dPです。"
	DrawMessage = "引き分けです。\nあなたの手札は [%s] です。合計値は %d。\nディーラーの手札は [%s] です。合計値は %d。\nあなたの所持金は %dPです。"

	InsuranceMessage    = "インシュランスしました。\nあなたの手札は [%s] です。合計値は %d。追加で %dP支払いました。\n あなたの所持金は %dPです。次の操作を行ってください"
	InsuranceWinMessage = "インシュランスに勝ちました。\nあなたの手札は [%s] です。合計値は %d。\nディーラーの手札は [%s] です。合計値は %d。\nあなたの所持金は %dPです。"

	BustMessage       = "バーストしました。\nあなたの手札は [%s] です。合計値は %d。\nディーラーの手札は [%s] です。合計値は %d。\nあなたの所持金は %dPです。"
	SplitMessage      = "スプリットしました。\nあなたの手札は [%s] [%s] です。合計値は %dと%dです。\n次の操作を行ってください"
	DoubleDownMessage = "ダブルダウンしました。\nあなたの手札は [%s] です。合計値は %d。"
	SurrenderMessage  = "サレンダーしました。"

	// Template Message
	BalanceMessage  = "あなたの所持金は %dPです。"
	NextStepMessage = "次の操作を行ってください"
	// Error Message
	InvalidMessage     = "無効な操作です。次の操作を行ってください"
	NoMoneyMessage     = "所持金が足りません。ゲームを終了します。"
	NoSplitMessage     = "スプリットできません。次の操作を行ってください"
	NoDoubleMessage    = "ダブルダウンできません。次の操作を行ってください"
	NoSurrenderMessage = "サレンダーできません。次の操作を行ってください"
	NoInsuranceMessage = "インシュランスできません。次の操作を行ってください"
	NoHitMessage       = "ヒットできません。次の操作を行ってください"
	NoStandMessage     = "スタンドできません。次の操作を行ってください"
	NoBetMessage       = "賭け金を設定してください。"
	NoStartMessage     = "ゲームを開始してください。"
)
