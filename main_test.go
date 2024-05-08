package main

import (
	"fmt"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/tool"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"regexp"
	"sync"
	"testing"
	"time"

	global "github.com/IUnlimit/telegram-bot-disrecall/internal"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/bot"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/conf"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/db"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/logger"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	log "github.com/sirupsen/logrus"
)

// go test -test.run TestMain
func TestMain(m *testing.M) {
	conf.Init()
	logger.Init()
	db.Init()
	//updateDB()
	fetchFile()
}

func fetchFile() {
	config := global.Config.TelegramBot
	b := bot.NewBasicTGBot(config.Token, config.Endpoint, config.Debug)
	log.Infof("Authorized on account %s(@%s)", b.API.Self.FirstName, b.API.Self.UserName)

	fileConfig := tgbotapi.FileConfig{FileID: "AgACAgUAAxkBAANqZiPugbNtfYPHafn1Ux5Sdsux2jEAAuewMRu-4clUF7enptP-awcBAAMCAAN5AAM0BA"}
	file, err := b.API.GetFile(fileConfig)
	if err != nil {
		log.Errorf("Fetch file info failed: %v", err)
		return
	}

	link := file.Link(b.API.Token)
	log.Info(link)
}

func updateDB() {
	config := global.Config.TelegramBot
	basic := bot.NewBasicTGBot(config.Token, config.Endpoint, config.Debug)
	log.Infof("Authorized on account %s(@%s)", basic.API.Self.FirstName, basic.API.Self.UserName)

	var wg sync.WaitGroup
	var fileModels []*model.FileModel
	// .Where("file_size = ?", 0)
	db.Instance.Where("file_type IN ('Video', 'Photo')").Find(&fileModels)
	log.Infof("Start %d tasks", len(fileModels))
	for _, fileModel := range fileModels {
		wg.Add(1)
		go func(fileModel *model.FileModel) {
			defer wg.Done()
			testDownload(fileModel, basic)
			// db.Instance.Model(fileModel).Updates(model.FileModel{FileSize: int64(file.FileSize), FilePath: file.FilePath})
		}(fileModel)
		break
	}
	wg.Wait()
}

func testDownload(fileModel *model.FileModel, b *bot.BasicTGBot) {
	fileConfig := tgbotapi.FileConfig{FileID: fileModel.FileID}
	// TODO 文件过大时 http 超时, 无法获取 Info ?
	file, err := b.API.GetFile(fileConfig)
	if err != nil {
		log.Errorf("Fetch file info failed: %v", err)
		return
	}

	//// 判断是否为本地服务器
	//if strings.HasPrefix(file.FilePath, "/") {
	//	// 将绝对路径转换为相对路径
	//	index := strings.Index(file.FilePath, b.API.Token)
	//	file.FilePath = fmt.Sprintf(".%s", file.FilePath[index+len(b.API.Token):])
	//}

	fileDirectURL := file.Link(b.API.Token)
	rootDir := global.Config.RootDir
	date := time.Now().Format("2006-01-02")
	fromUserID := fileModel.Json.Data().From.ID
	filePath := fmt.Sprintf("%s/%s/%d/%s", rootDir, date, fromUserID, file.FilePath)
	log.Debugf("FileDirectURL: %s", fileDirectURL)

	// 替换文件名
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
	log.Infof("文件成功下载到: %s", filePath)

	db.Instance.Model(&model.FileModel{}).Where("id = ?", fileModel.ID).Update("file_path", filePath)
}
