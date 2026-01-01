package main

import (
	"sync"
	"testing"

	global "github.com/IUnlimit/telegram-bot-disrecall/internal"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/bot"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/conf"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/db"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/logger"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

// go test -test.run TestMain
func TestMain(m *testing.M) {
	conf.Init()
	logger.Init()
	db.Init()
	updateDB()
	//fetchFile()
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
	db.Instance.Where("file_type IN ('Video', 'Photo') AND file_path = ''").Find(&fileModels)
	log.Infof("Start %d tasks", len(fileModels))
	for _, fileModel := range fileModels {
		wg.Add(1)
		go func(fileModel *model.FileModel) {
			defer wg.Done()
			testDownload(fileModel, basic)

		}(fileModel)
		//break
	}
	wg.Wait()
}

func testDownload(fileModel *model.FileModel, basic *bot.BasicTGBot) {
	_ = basic.DownloadFile(fileModel.FileID, fileModel.Json.Data(), func(filePath string, fileSize int64) {
		fileModel.FilePath = filePath
		db.Instance.Model(fileModel).Updates(model.FileModel{FileSize: fileSize, FilePath: filePath})
	})
}
