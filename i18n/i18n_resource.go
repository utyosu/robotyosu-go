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
		"ja": {
			"recruit":                "[%v] %v by %v (%v/%v)",
			"participants":           "    参加者: %s",
			"alternate_participants": "    補欠: %s",
			"no_recruitment":         "募集はありません",
			"not_join_because_full":  "[%v] は満員なので参加できませんでした(´・ω・`)ｼｮﾎﾞｰﾝ",
			"reserve_on_time":        "[%v] の開始時間になりました。あと%v人ｲﾅｲｶﾅ?(o・ω・)",
			"expired":                "[%v] は期限を過ぎたので終了しました。",
			"close_reserved":         "%v\n[%v] の開始時間になりました。いってらっしゃい！",
			"join":                   "%vさんが [%v] に参加しました。",
			"join_alternate":         "%vさんが [%v] に補欠で参加しました。",
			"leave":                  "%vさんが [%v] をキャンセルしました。",
			"gathered":               "%v\n[%v] のメンバーが集まったよ！(｀・ω・´)",
			"gathered_reserved":      "[%v] のメンバーが集まったよ！(｀・ω・´)予定時間になったらお知らせするね！",
			"closed":                 "%vさんが [%v] を終了しました。",
			"resurrection":           "最後に終了した募集 [%v] を再開しました。",
			"open_with_reserve":      "%vさんから [%v] を予定時間 %v で募集を受け付けました。",
			"open":                   "%vさんから [%v] を期限 %v で募集を受け付けました。",
			"too_long_title":         "募集メッセージが長すぎます。",
			"capacity_less":          "募集人数は1人以上にする必要があります。",
			"capacity_over":          "募集人数が多すぎます。",
			"reserve_limit_over":     "予定日時が遠すぎます。",
			"multiple_number":        "数字が複数あるので分かりませんでした(´・ω・`)ｼｮﾎﾞｰﾝ",
			"twitter_recruit":        "%v\n%v by %v (%v/%v)",
			"twitter_members":        "参加者: %v",
			"twitter_close":          "%v\nこの募集は終了しました。",
			"error":                  "エラーが発生しました。何度も発生する場合は開発者にお問い合わせ下さい。",
			"twitter_error":          "ツイートの投稿に失敗しました。何度も発生する場合は開発者にお問い合わせ下さい。",
		},
		"en": {
			"recruit":                "[%v] %v by %v (%v/%v)",
			"participants":           "    Members: %s",
			"alternate_participants": "    Alternate: %s",
			"no_recruitment":         "No recruitments.",
			"not_join_because_full":  "You cannot join [%v] because full. :_(",
			"reserve_on_time":        "It's time to start [%v]. We are looking for %v more. ('')>",
			"expired":                "[%v] is expired.",
			"close_reserved":         "%v\nIt's time to start [%v]. Good luck!",
			"join":                   "%v joined [%v].",
			"join_alternate":         "%v alternate joined [%v].",
			"leave":                  "%v leave [%v].",
			"gathered":               "%v\n[%v] is gathered. :)",
			"gathered_reserved":      "[%v] is gathered. Notify member when it's time. XD",
			"closed":                 "%v closed [%v].",
			"resurrection":           "The last recruitment [%v], has been resumed.",
			"open_with_reserve":      "%v open [%v], reserved at %v.",
			"open":                   "%v open [%v], expire at %v.",
			"too_long_title":         "Too long recruitment subject.",
			"capacity_less":          "The number of applicants cannot be less than one.",
			"capacity_over":          "Capacity over.",
			"reserve_limit_over":     "Too long scheduled time.",
			"multiple_number":        "There are multiple numbers.",
			"twitter_recruit":        "%v\n%v by %v (%v/%v)",
			"twitter_members":        "members: %v",
			"twitter_close":          "%v\nThis recruitment is closed.",
			"error":                  "Cause error.",
			"twitter_error":          "Post tweet error.",
		},
	}

	helpBasicCommands = map[string][]commandSet{
		"ja": {
			{".rt enable", "ロボちょすBOTの有効化"},
			{".rt disable", "ロボちょすBOTの無効化"},
			{".rt help", "設定変更コマンドの参照"},
			{"", "有効化後に使えるコマンド"},
			{".rt language", "言語の参照"},
			{".rt language ${LANGUAGE}", "言語の変更"},
			{".rt timezone", "タイムゾーンの参照"},
			{".rt timezone ${TIMEZONE}", "タイムゾーンの変更"},
			{".rt reserve_limit_time", "募集の予定日時の制限の参照 (指定した秒数より先の募集ができなくなる。0を指定すると制限なし)"},
			{".rt reserve_limit_time ${TIME}", "募集の予定日時の制限の変更"},
			{".rt expire_duration", "日時指定をせずに募集したときの締め切りまでの時間（秒数）"},
			{".rt expire_duration ${TIME}", "expire_durationの変更"},
			{".rt expire_duration_for_reserve", "日時指定をして募集したときの締め切りまでの時間（秒数）"},
			{".rt expire_duration_for_reserve ${TIME}", "expire_duration_for_reserveの変更"},
		},
		"en": {
			{".rt enable", "Enable robotyosu bot."},
			{".rt disable", "Disable robotyosu bot."},
			{".rt help", "Show setting commands."},
			{"", "Available after enables."},
			{".rt language", "Show language."},
			{".rt language ${LANGUAGE}", "Change language."},
			{".rt timezone", "Show timezone."},
			{".rt timezone ${TIMEZONE}", "Change timezone."},
			{".rt reserve_limit_time", "Show reserve limit time. Recruitment after the specified number of seconds will not be possible. 0 is unlimited."},
			{".rt reserve_limit_time ${TIME}", "Change reserve limit time."},
			{".rt expire_duration", "Expire duration for non-reserved recruitment. (sec)"},
			{".rt expire_duration ${TIME}", "Change expire duration."},
			{".rt expire_duration_for_reserve", "Expire duration for reserved recruitment. (sec)"},
			{".rt expire_duration_for_reserve ${TIME}", "Change expire duration for reserve."},
		},
	}

	helpRecruitmentCommands = map[string][]commandSet{
		"ja": {
			{"募集内容@<数字>", "募集の開始 (例「ゲームしましょう@3」)"},
			{"<数字>参加", "募集に参加 (例「1参加」)"},
			{"<数字>キャンセル", "参加キャンセル (例「1キャンセル」)"},
			{"<数字>しめ", "募集の終了 (例「1しめ」)"},
			{"復活", "最後に終了した募集を再開"},
			{"案件", "最新の募集状態を表示"},
		},
		"en": {
			{"Recruitment contents@<number>", "Start recruitment (ex. Play games@3)"},
			{"<number>join", "Join recruitment (ex. 1 join)"},
			{"<number>cancel", "Cancel participation (ex. 1 cancel)"},
			{"<number>close", "Close recruitment (ex. 1 close)"},
			{"resume", "Resume the last closed recruitment"},
			{"list", "Show recruitments"},
		},
	}
)
