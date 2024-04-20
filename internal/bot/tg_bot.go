package bot

import (
	"fmt"
	"os"
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

func (b *BasicTGBot) SendMessage(text string, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	b.API.Send(msg)
}

func (b *BasicTGBot) DownloadFile(message *tgbotapi.Message, fileID string) {
	// 获取文件信息
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	file, err := b.API.GetFile(fileConfig)
	if err != nil {
		log.Error("文件信息获取失败:", err)
		return
	}

	// 判断是否为本地服务器
	if strings.HasPrefix(file.FilePath, "/") {
		b.SendMessage("文件被本地服务器成功保存到: "+file.FilePath, message)
		return
	}

	fileDirectURL := file.Link(b.API.Token)
	rootDir := global.Config.RootDir
	date := time.Now().Format("2006-01-02")
	fromUserID := message.From.ID
	filePath := fmt.Sprintf("%s/%s/%d/%s", rootDir, date, fromUserID, file.FilePath)
	log.Debugf("FileDirectURL: %s", fileDirectURL)

	// 下载文件
	filePath, err = tool.DownloadFile(fileDirectURL, filePath)
	if err != nil {
		log.Println("文件 "+fileDirectURL+" 下载失败:", err)
		_ = os.Remove(filePath)
		return
	}

	log.Println("文件" + filePath + "成功下载!")
	b.SendMessage("文件成功下载到 "+filePath, message)
}
