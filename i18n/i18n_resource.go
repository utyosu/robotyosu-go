package i18n

type commandSet struct {
	command     string
	description string
}

const (
	DefaultLanguage = "ja"
)

var (
	Languages = []string{"ja", "en"}

	commonDictionary = map[string]string{
		"enable":  "募集機能を有効にしました。使い方を見るには `.rt help` と入力して下さい。\nEnable recruitment. Type `.rt help` to see how to use it.\nType `.rt language en` to change the language to English.",
		"disable": "募集機能を無効にしました。\nDisable recruitment.",
	}

	languageDictionary = map[string]map[string]string{
		"ja": map[string]string{
			"recruit":               "[%v] %v by %v",
			"participants":          "    参加者: %s",
			"no_recruitment":        "募集はありません",
			"not_join_because_full": "[%v] は満員なので参加できませんでした(´・ω・`)ｼｮﾎﾞｰﾝ",
			"reserve_on_time":       "[%v] の開始時間になりました。あと%v人ｲﾅｲｶﾅ?(o・ω・)",
			"expired":               "[%v] は期限を過ぎたので終了しました。",
			"close_reserved":        "%v\n[%v] の開始時間になりました。いってらっしゃい！",
			"join":                  "%vさんが [%v] に参加しました。",
			"leave":                 "%vさんが [%v] をキャンセルしました。",
			"gathered":              "%v\n[%v] のメンバーが集まったよ！(｀・ω・´)",
			"gathered_reserved":     "[%v] のメンバーが集まったよ！(｀・ω・´)予定時間になったらお知らせするね！",
			"closed":                "%vさんが [%v] を終了しました。",
			"resurrection":          "最後に終了した募集 [%v] を再開しました。",
			"open_with_reserve":     "%vさんから [%v] を予定時間 %v で募集を受け付けました。",
			"open":                  "%vさんから [%v] を期限 %v で募集を受け付けました。",
			"too_long_title":        "募集メッセージが長すぎます。",
			"capacity_less":         "募集人数は1人以上にする必要があります。",
			"capacity_over":         "募集人数が多すぎます。",
			"twitter_recruit":       "%v\n%v by %v",
			"twitter_members":       "参加者: %v",
			"twitter_close":         "%v\nこの募集は終了しました。",
			"error":                 "エラーが発生しました。何度も発生する場合は開発者にお問い合わせ下さい。",
		},
		"en": map[string]string{
			"recruit":               "[%v] %v by %v",
			"participants":          "    Members: %s",
			"no_recruitment":        "No recruitments.",
			"not_join_because_full": "You cannot join [%v] because full. :_(",
			"reserve_on_time":       "It's time to start [%v]. We are looking for %v more. ('')>",
			"expired":               "[%v] is expired.",
			"close_reserved":        "%v\nIt's time to start [%v]. Good luck!",
			"join":                  "%v joined [%v].",
			"leave":                 "%v leave [%v].",
			"gathered":              "%v\n[%v] is gathered. :)",
			"gathered_reserved":     "[%v] is gathered. Notify member when it's time. XD",
			"closed":                "%v closed [%v].",
			"resurrection":          "The last recruitment [%v], has been resumed.",
			"open_with_reserve":     "%v open [%v], reserved at %v.",
			"open":                  "%v open [%v], expire at %v.",
			"too_long_title":        "Too long recruitment subject.",
			"capacity_less":         "The number of applicants cannot be less than one.",
			"capacity_over":         "Capacity over.",
			"twitter_recruit":       "%v\n%v by %v",
			"twitter_members":       "members: %v",
			"twitter_close":         "%v\nThis recruitment is closed.",
			"error":                 "Cause error.",
		},
	}

	helpBasicCommands = map[string][]commandSet{
		"ja": []commandSet{
			{".rt enable", "ロボちょすBOTの有効化"},
			{".rt disable", "ロボちょすBOTの無効化"},
			{".rt help", "設定変更コマンドの参照"},
			{"", "有効化後に使えるコマンド"},
			{".rt language", "言語の参照"},
			{".rt language ${LANGUAGE}", "言語の変更"},
			{".rt timezone", "タイムゾーンの参照"},
			{".rt timezone ${TIMEZONE}", "タイムゾーンの変更"},
			{"使い方", "募集機能の使い方"},
		},
		"en": []commandSet{
			{".rt enable", "Enable robotyosu bot"},
			{".rt disable", "Disable robotyosu bot"},
			{".rt help", "Show setting commands."},
			{"", "Available after enables"},
			{".rt language", "Show language"},
			{".rt language ${LANGUAGE}", "Change language"},
			{".rt timezone", "Show timezone"},
			{".rt timezone ${TIMEZONE}", "Change timezone"},
			{"help", "How to use recruitment"},
		},
	}

	helpRecruitmentCommands = map[string][]commandSet{
		"ja": []commandSet{
			{"募集内容@<数字>", "募集の開始 (例「ゲームしましょう@3」)"},
			{"<数字>参加", "募集に参加 (例「1参加」)"},
			{"<数字>キャンセル", "参加キャンセル (例「1キャンセル」)"},
			{"<数字>しめ", "募集の終了 (例「1しめ」)"},
			{"復活", "最後に終了した募集を再開"},
			{"案件", "最新の募集状態を表示"},
		},
		"en": []commandSet{
			{"Recruitment contents@<number>", "Start recruitment (ex. Play games@3)"},
			{"<number>join", "Join recruitment (ex. 1 join)"},
			{"<number>cancel", "Cancel participation (ex. 1 cancel)"},
			{"<number>close", "Close recruitment (ex. 1 close)"},
			{"resume", "Resume the last closed recruitment"},
			{"list", "Show recruitments"},
		},
	}
)
