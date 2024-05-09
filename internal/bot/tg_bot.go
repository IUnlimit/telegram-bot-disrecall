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
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	// TODO 文件过大时 http 超时, 无法获取 Info ?
	file, err := b.API.GetFile(fileConfig)
	if err != nil {
		log.Errorf("Fetch file info failed: %v", err)
		return
	}

	rootDir := global.Config.RootDir
	date := time.Now().Format("2006-01-02")
	fromUserID := message.From.ID

	// if local bot api server
	if strings.HasPrefix(file.FilePath, "/") {
		log.Debugf("Received file from local server: %s", file.FilePath)
		selfAPIConfig := global.Config.SelfAPI
		// 将映射路径转换为真实绝对路径
		filePath := strings.Replace(file.FilePath, selfAPIConfig.RootDir, selfAPIConfig.RealRootDir, 1)
		// 构造相对路径, 防呆
		copyRelativePath := strings.Replace(file.FilePath, selfAPIConfig.RootDir, "./", 1)
		copyRelativePath = strings.Replace(copyRelativePath, b.API.Token, "", 1)
		copyPath := generateFileName(fmt.Sprintf("%s/%s/%d/%s", rootDir, date, fromUserID, copyRelativePath), file.FileUniqueID)

		err = tool.CopyFile(filePath, copyPath)
		if err != nil {
			log.Errorf("Copy file from '%s' to '%s' failed: %v", filePath, copyPath, err)
			return
		}

		log.Infof("File successfully copied to '%s'", copyPath)
		b.SendMessage(fmt.Sprintf("文件成功拷贝到: %s", copyPath), message)
		callback(copyPath, int64(file.FileSize))
		return
	}

	fileDirectURL := file.Link(b.API.Token)
	filePath := generateFileName(fmt.Sprintf("%s/%s/%d/%s", rootDir, date, fromUserID, file.FilePath), file.FileUniqueID)
	log.Debugf("FileDirectURL: %s", fileDirectURL)

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

func generateFileName(path string, fileUID string) string {
	// 定义正则表达式，匹配文件名
	re := regexp.MustCompile(`(.*)(file_\d+)(\..*)?$`)
	if re.NumSubexp() != 3 {
		log.Errorf("Can't match path '%s' with regexp: %s", path, re.String())
		return path
	}
	// fileName -> fileUID
	result := re.ReplaceAllString(path, "${1}"+fileUID+"$3")
	// 没有后缀
	if re.FindStringSubmatch(path)[3] == "" {
		if strings.Contains(path, "photos/") {
			result += ".jpg"
		} else if strings.Contains(path, "videos/") {
			result += ".mp4"
		} else {
			log.Warnf("Unknown file type with path: %s", path)
		}
	}
	return result
}
