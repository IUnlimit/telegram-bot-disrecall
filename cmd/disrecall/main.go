package disrecall

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"gorm.io/datatypes"

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
			if up.UpdateID >= u.Offset {
				u.Offset = up.UpdateID + 1
				go func(update *tgbotapi.Update, bot *bot.BasicTGBot) {
					handle(&current, size, bot, update)
				}(&up, basicBot)
			}
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
	unique := message.MediaGroupID
	if unique == "" {
		unique = strconv.Itoa(message.ForwardDate)
	}

	filePairs := findFileIDWithType(message)
	if len(filePairs) == 0 {
		basic.SendMessage("当前消息类型不支持储存!", message)
		return
	}

	for _, filePair := range filePairs {
		go func(message *tgbotapi.Message, filePair *model.Pair[string, model.FileType]) {
			fileModel := &model.FileModel{
				Model: gorm.Model{
					CreatedAt: time.Unix(int64(message.Date), 0),
				},
				MessageID:    int64(message.MessageID),
				UserID:       message.From.ID,
				ForwardDate:  int64(message.ForwardDate),
				MediaGroupID: unique,
				UserName:     message.From.UserName,
				Text:         message.Text,
				Json:         datatypes.NewJSONType(message),
			}

			var err error
			fileID := filePair.Key
			fileType := filePair.Value

			if fileType == model.Text {
				fileModel.FileType = fileType
				fileModel.Text = message.Text
			} else {
				for i := 0; i < 3; i++ {
					err = basic.DownloadFile(fileID, message, func(filePath string, fileSize int64) {
						fileModel.FilePath = filePath
						fileModel.FileType = fileType
						fileModel.FileSize = fileSize
						fileModel.FileID = fileID
					})
					if err == nil {
						break
					}
					log.Errorf("Failed to download file with retry-%d: %v", i, err)
				}
			}

			if err != nil {
				basic.SendMessage(fmt.Sprintf("Batch(%d): 消息下载失败 %v", message.ForwardDate, err), message)
				return
			}

			err = db.InsertFile(fileModel)
			if err != nil {
				log.Errorf("Database insert error: %v", err)
				basic.SendMessage(fmt.Sprintf("Batch(%d): 数据库录入失败 %v", message.ForwardDate, err), message)
				return
			}

			atomic.AddInt32(current, 1)
			basic.SendMessage(fmt.Sprintf("Batch(%d): 消息已录入数据库 [%d/%d]", message.ForwardDate, *current, size), message)
		}(message, filePair)
	}
}

// TODO support multiple photos
func findFileIDWithType(message *tgbotapi.Message) []*model.Pair[string, model.FileType] {
	filePairs := make([]*model.Pair[string, model.FileType], 0)

	if message.Photo != nil && len(message.Photo) > 0 {
		// the same photo with diff size
		maxPhoto := message.Photo[0]
		for _, photo := range message.Photo {
			if photo.FileSize > maxPhoto.FileSize {
				maxPhoto = photo
			}
		}
		filePairs = append(filePairs, model.NewPair(maxPhoto.FileID, model.Photo))
	} else if message.Voice != nil {
		filePairs = append(filePairs, model.NewPair(message.Voice.FileID, model.Voice))
	} else if message.Video != nil {
		filePairs = append(filePairs, model.NewPair(message.Video.FileID, model.Video))
	} else if message.Document != nil {
		filePairs = append(filePairs, model.NewPair(message.Document.FileID, model.Document))
	} else if message.Text != "" {
		filePairs = append(filePairs, model.NewPair("", model.Text))
	}
	return filePairs
}
