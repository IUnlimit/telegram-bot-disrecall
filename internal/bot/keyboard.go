package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var mainReplyKeyboard tgbotapi.InlineKeyboardMarkup

// 查詢指令内联键盘
func generateQueryInlineKeyboard(tag string, current int, total int) *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("%s#%d", tag, current-1)),
			tgbotapi.NewInlineKeyboardButtonData("⌘", HelpCommand),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d⁄%d", current, total), fmt.Sprintf(">%s#1-%d", tag, total)),
			tgbotapi.NewInlineKeyboardButtonData("➡️", fmt.Sprintf("%s#%d", tag, current+1)),
		),
	)
	return &markup
}
