package bot

import (
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init() {
	InputChanMap = make(map[int64]chan string)

	supportTypes = map[string]func(tgbotapi.RequestFileData, *model.FileModel) interface{}{
		string(model.Text): func(_ tgbotapi.RequestFileData, fileModel *model.FileModel) interface{} {
			text := tgbotapi.NewMessage(0, fileModel.Text)
			// TODO fix message.Entities, media.CaptionEntities 无法生效
			return &text
		},
		string(model.Photo): func(file tgbotapi.RequestFileData, fileModel *model.FileModel) interface{} {
			media := tgbotapi.NewInputMediaPhoto(file)
			message := fileModel.Json.Data()
			media.Caption = message.Caption
			return &media
		},
		string(model.Voice): func(file tgbotapi.RequestFileData, fileModel *model.FileModel) interface{} {
			media := tgbotapi.NewInputMediaAudio(file)
			message := fileModel.Json.Data()
			media.Caption = message.Caption
			return &media
		},
		string(model.Video): func(file tgbotapi.RequestFileData, fileModel *model.FileModel) interface{} {
			media := tgbotapi.NewInputMediaVideo(file)
			message := fileModel.Json.Data()
			media.Caption = message.Caption
			return &media
		},
		string(model.Document): func(file tgbotapi.RequestFileData, fileModel *model.FileModel) interface{} {
			media := tgbotapi.NewInputMediaDocument(file)
			message := fileModel.Json.Data()
			media.Caption = message.Caption
			return &media
		},
	}

	cmdFuncMap = map[string]func(context *CommandContext){
		StartCommand:  onStart,
		HelpCommand:   onHelp,
		StaticCommand: onStatic,
		ListTextsCommand: func(c *CommandContext) {
			onListMedia(model.Text, c)
		},
		ListPhotosCommand: func(c *CommandContext) {
			onListMedia(model.Photo, c)
		},
		ListVoicesCommand: func(c *CommandContext) {
			onListMedia(model.Voice, c)
		},
		ListVideosCommand: func(c *CommandContext) {
			onListMedia(model.Video, c)
		},
		ListDocumentsCommand: func(c *CommandContext) {
			onListMedia(model.Document, c)
		},
	}

	// TODO 支持删除记录
	mainReplyKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("查看文本存档", ListTextsCommand),
			tgbotapi.NewInlineKeyboardButtonData("查看图片存档", ListPhotosCommand),
			tgbotapi.NewInlineKeyboardButtonData("查看视频存档", ListVideosCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("查看音频存档", ListVoicesCommand),
			// TODO
			tgbotapi.NewInlineKeyboardButtonData("查看消息存档", "/multiples"),
			tgbotapi.NewInlineKeyboardButtonData("查看文件存档", ListDocumentsCommand),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("统计数据", StaticCommand),
		),
	)
}
