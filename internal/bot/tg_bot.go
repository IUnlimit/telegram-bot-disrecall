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

func (b *BasicTGBot) SendMessage(text string, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyToMessageID = message.MessageID
	b.API.Send(msg)
}

func (b *BasicTGBot) DownloadFile(fileID string, message *tgbotapi.Message, callback func(filePath string)) {
	// 获取文件信息
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	file, err := b.API.GetFile(fileConfig)
	if err != nil {
		log.Errorf("Fetch file info failed: %v", err)
		return
	}

	// 判断是否为本地服务器
	if strings.HasPrefix(file.FilePath, "/") {
		// 绝对路径要去 token
		b.SendMessage(fmt.Sprintf("文件被本地服务器成功保存到: %s", replaceToken(file.FilePath)), message)
		callback(file.FilePath)
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
		log.Errorf("File %s download failed: %v", fileDirectURL, err)
		_ = os.Remove(filePath)
		return
	}

	log.Infof("File %s download success", filePath)
	b.SendMessage(fmt.Sprintf("文件成功下载到: %s", filePath), message)
	callback(filePath)
}

func replaceToken(path string) string {
	re := regexp.MustCompile(`/\d+:[^/]+/`)
	return re.ReplaceAllString(path, "/<token>/")
}
