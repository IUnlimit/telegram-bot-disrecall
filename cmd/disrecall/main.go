package disrecall

import (
	"fmt"
	"gorm.io/datatypes"
	"sync/atomic"
	"time"

	"github.com/kylelemons/godebug/pretty"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	global "github.com/IUnlimit/telegram-bot-disrecall/internal"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/bot"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/db"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Listen 不支持 Channel 消息, 仅私聊可用
func Listen() {
	config := global.Config.TelegramBot
	basicBot := bot.NewBasicTGBot(config.Token, config.Endpoint, config.Debug)
	log.Infof("Authorized on account %s(@%s)", basicBot.API.Self.FirstName, basicBot.API.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for {
		updates, err := basicBot.API.GetUpdates(u)
		if err != nil {
			log.Errorf("Failed to get updates, retrying in 3 seconds...: %v", err)
			time.Sleep(time.Second * 3)
			continue
		}

		current := int32(0)
		size := len(updates)
		for _, up := range updates {
			go func(update *tgbotapi.Update, config *tgbotapi.UpdateConfig) {
				if update.UpdateID >= config.Offset {
					config.Offset = update.UpdateID + 1
					handle(&current, size, basicBot, update)
				}
			}(&up, &u)
		}
	}
}

func handle(current *int32, size int, basic *bot.BasicTGBot, update *tgbotapi.Update) {
	if update.CallbackQuery != nil {
		bot.OnCommandCallback(update.CallbackQuery, basic)
		return
	}

	// ignore any non-Message updates
	if update.Message == nil {
		return
	}

	// input capture
	if channel, ok := bot.InputChanMap[update.Message.Chat.ID]; ok {
		channel <- update.Message.Text
		return
	}

	// user auth
	if !bot.UserAuth(update.Message.Chat.ID, basic) {
		return
	}

	if update.Message.IsCommand() {
		bot.OnCommand(update.Message, basic)
		return
	}

	log.Debugf("[%s] %v", update.Message.From.UserName, pretty.Sprint(update.Message))

	if update.Message.ForwardDate == 0 {
		basic.SendMessage("未获取到转发消息数据!", update.Message)
		return
	}

	message := update.Message
	fileModel := &model.FileModel{
		Model: gorm.Model{
			CreatedAt: time.Unix(int64(message.Date), 0),
		},
		MessageID:   int64(message.MessageID),
		UserID:      message.From.ID,
		ForwardDate: int64(message.ForwardDate),
		UserName:    message.From.UserName,
		Text:        message.Text,
		Json:        datatypes.NewJSONType(message),
	}

	fileID, fileType := findFileIDWithType(message)
	if fileType == "" {
		basic.SendMessage("当前消息类型不支持储存!", message)
		return
	}

	go func(message *tgbotapi.Message, fileModel *model.FileModel) {
		if fileType == model.Text {
			fileModel.FileType = fileType
			fileModel.Text = message.Text
		} else {
			basic.DownloadFile(fileID, message, func(filePath string, fileSize int64) {
				fileModel.FilePath = filePath
				fileModel.FileType = fileType
				fileModel.FileSize = fileSize
				fileModel.FileID = fileID
			})
		}

		err := db.InsertFile(fileModel)
		if err != nil {
			log.Errorf("Database insert error: %v", err)
			return
		}

		atomic.AddInt32(current, 1)
		basic.SendMessage(fmt.Sprintf("Batch(%d): 消息已录入数据库 [%d/%d]", message.ForwardDate, *current, size), message)
	}(message, fileModel)
}

func findFileIDWithType(message *tgbotapi.Message) (string, model.FileType) {
	if message.Photo != nil {
		photo := (message.Photo)[len(message.Photo)-1]
		return photo.FileID, model.Photo
	}

	if message.Voice != nil {
		return message.Voice.FileID, model.Voice
	} else if message.Video != nil {
		return message.Video.FileID, model.Video
	} else if message.Document != nil {
		return message.Document.FileID, model.Document
	} else if message.Text != "" {
		return "", model.Text
	}
	return "", ""
}
