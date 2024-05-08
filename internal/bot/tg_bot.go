package bot

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	global "github.com/IUnlimit/telegram-bot-disrecall/internal"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/tool"
	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BasicTGBot struct {
	API *tgbotapi.BotAPI
}

func NewBasicTGBot(token string, endpoint string, debug bool) *BasicTGBot {
	bot, err := tgbotapi.NewBotAPIWithAPIEndpoint(token, endpoint)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = debug

	return &BasicTGBot{
		API: bot,
	}
}

func (b *BasicTGBot) Send(chattable tgbotapi.Chattable) {
	_, err := b.API.Send(chattable)
	if err != nil {
		log.Info(err)
	}
}

func (b *BasicTGBot) SendMessage(text string, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	b.Send(msg)
}

func (b *BasicTGBot) SendChatMessage(text string, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.Send(msg)
}

func (b *BasicTGBot) SendCallback(text string, callback *tgbotapi.CallbackQuery) {
	msg := tgbotapi.NewCallback(callback.ID, text)
	b.Send(msg)
}

func (b *BasicTGBot) DownloadFile(fileID string, message *tgbotapi.Message, callback func(filePath string, fileSize int64)) {
	// 获取文件信息
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	// TODO 文件过大时 http 超时, 无法获取 Info ?
	file, err := b.API.GetFile(fileConfig)
	if err != nil {
		log.Errorf("Fetch file info failed: %v", err)
		return
	}

	// 判断是否为本地服务器
	if strings.HasPrefix(file.FilePath, "/") {
		// 将绝对路径转换为相对路径
		index := strings.Index(file.FilePath, b.API.Token)
		file.FilePath = fmt.Sprintf(".%s", file.FilePath[index+len(b.API.Token):])
		// TODO move file if local mode
		return
	}

	fileDirectURL := file.Link(b.API.Token)
	rootDir := global.Config.RootDir
	date := time.Now().Format("2006-01-02")
	fromUserID := message.From.ID
	filePath := fmt.Sprintf("%s/%s/%d/%s", rootDir, date, fromUserID, file.FilePath)
	log.Debugf("FileDirectURL: %s", fileDirectURL)

	// 替换文件名 file_11.jpg -> <id>.jpg
	re := regexp.MustCompile(`([^/]+)\.([^.]+)$`)
	filePath = re.ReplaceAllString(filePath, file.FileUniqueID+".$2")

	// 下载文件
	filePath, err = tool.DownloadFile(fileDirectURL, filePath)
	if err != nil {
		log.Errorf("File %s download failed: %v", fileDirectURL, err)
		_ = os.Remove(filePath)
		return
	}

	log.Infof("File successfully download to '%s'", filePath)
	b.SendMessage(fmt.Sprintf("文件成功下载到: %s", filePath), message)
	callback(filePath, int64(file.FileSize))
}
