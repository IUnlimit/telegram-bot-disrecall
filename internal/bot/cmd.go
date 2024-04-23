package bot

import (
	"fmt"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/cache"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	StartCommand         = "/start"
	HelpCommand          = "/help"
	StaticCommand        = "/static"
	ListTextsCommand     = "/texts"
	ListPhotosCommand    = "/photos"
	ListVoicesCommand    = "/voices"
	ListVideosCommand    = "/videos"
	ListDocumentsCommand = "/docs"
)

var cmdFuncMap map[string]func(*CommandContext)

type CommandContext struct {
	Response *tgbotapi.MessageConfig
	Message  *tgbotapi.Message
	Basic    *BasicTGBot
}

func OnCommand(message *tgbotapi.Message, basic *BasicTGBot) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	context := &CommandContext{
		Response: &msg,
		Message:  message,
		Basic:    basic,
	}

	if cmdFunc, ok := cmdFuncMap[message.Text]; ok {
		cmdFunc(context)
	} else {
		msg.Text = "未知的命令: " + message.Text
	}

	if _, err := basic.API.Send(msg); err != nil {
		log.Error(err)
	}
}

func onStart(context *CommandContext) {
	context.Response.Text = "您好，我是防撤回机器人，您可将需要保存的 文本/图片/语音/视频/文件 转发给我，我会自动将文件存档到本地。即使被撤回，存档文件也可重新被查阅"
	context.Response.ReplyMarkup = mainReplyKeyboard
}

func onHelp(context *CommandContext) {
	context.Response.Text = "帮助菜单"
	context.Response.ReplyMarkup = mainReplyKeyboard
}

func onStatic(context *CommandContext) {
	staticMap := cache.GetUserStatic(context.Message.Chat.ID)
	var builder strings.Builder
	builder.WriteString("用户 @")
	builder.WriteString(context.Message.Chat.UserName)
	builder.WriteString("\n")
	for fileType, static := range staticMap {
		builder.WriteString("- ")
		builder.WriteString(string(fileType))
		builder.WriteString(": 共计 ")
		builder.WriteString(fmt.Sprintf("%.2f MB", static.MB))
		builder.WriteString(fmt.Sprintf("(%d 条)\n", static.Rows))
	}
	context.Response.Text = builder.String()
	context.Response.ParseMode = "Markdown"
	context.Response.ReplyMarkup = mainReplyKeyboard
}

func onListMedia(fileType model.FileType, context *CommandContext) {
	// 如果是 callback 进来，机器人会引用自身消息回复，此处 userID != FromID
	ListRecord(1, string(fileType), context.Message.Chat.ID, context.Message.Chat.ID, context.Basic)
}
