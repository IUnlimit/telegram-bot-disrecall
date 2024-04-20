package disrecall

import (
	"encoding/json"
	"sync"
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

// 不支持 Channel 消息, 仅私聊可用
func Listen() {
	config := global.Config.TelegramBot
	basicBot := bot.NewBasicTGBot(config.Token, config.Endpoint, config.Debug)
	log.Infof("Authorized on account %s(@%s)", basicBot.API.Self.FirstName, basicBot.API.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := basicBot.API.GetUpdatesChan(u)
	for update := range updates {
		go handle(basicBot, &update)
	}
}

func handle(bot *bot.BasicTGBot, update *tgbotapi.Update) {
	// ignore any non-Message updates
	if update.Message == nil {
		return
	}
	log.Debugf("[%s] %v", update.Message.From.UserName, pretty.Sprint(update.Message))

	if update.Message.Text == "/start" {
		bot.SendMessage("您好，我是防撤回机器人，您可将需要保存的 文本/图片/语音/视频/文件 转发至此机器人，机器人将会自动将文件存档到本地服务器。若原消息被撤回，则机器人会将存档文件重新上传至该聊天", update.Message)
		return
	}

	if update.Message.ForwardDate == 0 {
		bot.SendMessage("未获取到转发消息数据!", update.Message)
		return
	}

	var wg sync.WaitGroup
	message := update.Message
	bytes, err := json.Marshal(&message)
	if err != nil {
		log.Errorf("Message(id: %d) marshall failed: %v", message.MessageID, err)
		return
	}

	fileModel := &model.FileModel{
		Model: gorm.Model{
			CreatedAt: time.Unix(int64(message.Date), 0),
		},
		MessageID:       int64(message.MessageID),
		ForwardUserID:   message.From.ID,
		ForwardDate:     int64(message.ForwardDate),
		ForwardUserName: message.From.UserName,
		Text:            message.Text,
		Caption:         message.Caption,
		Json:            string(bytes),
	}

	// 一条消息只包含一个副文本类型, 所以其实不用加协程 & 可以直接用 else if
	if message.Photo != nil {
		wg.Add(1)
		go func(message *tgbotapi.Message, fileModel *model.FileModel) {
			defer wg.Done()
			// 获取最后一张 (最清晰) 图片的 file_id
			photo := (message.Photo)[len(message.Photo)-1]
			fileID := photo.FileID
			bot.DownloadFile(fileID, message, func(filePath string) {
				fileModel.FilePath = filePath
				fileModel.FileType = model.Photo
			})
		}(message, fileModel)
	}

	if message.Voice != nil {
		wg.Add(1)
		go func(message *tgbotapi.Message, fileModel *model.FileModel) {
			defer wg.Done()
			fileID := message.Voice.FileID
			bot.DownloadFile(fileID, message, func(filePath string) {
				fileModel.FilePath = filePath
				fileModel.FileType = model.Voice
			})
		}(message, fileModel)
	}

	if message.Video != nil {
		wg.Add(1)
		go func(message *tgbotapi.Message, fileModel *model.FileModel) {
			defer wg.Done()
			fileID := message.Video.FileID
			bot.DownloadFile(fileID, message, func(filePath string) {
				fileModel.FilePath = filePath
				fileModel.FileType = model.Video
			})
		}(message, fileModel)
	}

	if message.Document != nil {
		wg.Add(1)
		go func(message *tgbotapi.Message, fileModel *model.FileModel) {
			defer wg.Done()
			fileID := message.Document.FileID
			bot.DownloadFile(fileID, message, func(filePath string) {
				fileModel.FilePath = filePath
				fileModel.FileType = model.Document
			})
		}(message, fileModel)
	}

	wg.Wait()

	err = db.Insert(fileModel)
	if err != nil {
		log.Errorf("Database insert error: %v", err)
	}
}
