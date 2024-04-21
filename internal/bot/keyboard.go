package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var replyKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/photos"),
		tgbotapi.NewKeyboardButton("/videos"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/texts"),
		tgbotapi.NewKeyboardButton("/voices"),
		tgbotapi.NewKeyboardButton("/docs"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/help"),
	),
)

// /photos 指令内联键盘
// current: 0, 1, 2, 3
func generatePhotosInlineKeyboard(tag string, current int, total int) *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("%s#%d", tag, current-1)),
			// TODO 点击捕获用户输入: 请输入你要跳转到的页面 (1-total)
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d⁄%d", current, total), NULL_CALLBACK),
			tgbotapi.NewInlineKeyboardButtonData("➡️", fmt.Sprintf("%s#%d", tag, current+1)),
		),
	)
	return &markup
}
