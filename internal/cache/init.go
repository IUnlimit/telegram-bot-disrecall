package cache

import (
	"github.com/IUnlimit/telegram-bot-disrecall/internal/db"
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// userID: FileModel[]
var cacheMap map[int64][]*model.FileModel

func Init() {
	timer := cron.New(cron.WithSeconds())
	updateLocalCache()
	_, err := timer.AddFunc("@every 5m", func() {
		// TODO 没有插入记录的不刷新
		updateLocalCache()
	})
	if err != nil {
		log.Fatalf("Creat cron task error: %v", err)
	}
	timer.Start()
}

func updateLocalCache() {
	newCacheMap := make(map[int64][]*model.FileModel, 0)
	fileModels, err := db.QueryAllFile()
	if err != nil {
		log.Fatalf("Initial cache data failed: %v", err)
	}

	var rows int
	for _, fileModel := range fileModels {
		if _, ok := newCacheMap[fileModel.UserID]; !ok {
			newCacheMap[fileModel.UserID] = make([]*model.FileModel, 0)
		}
		if !fileModel.IsValid() {
			continue
		}
		newCacheMap[fileModel.UserID] = append(newCacheMap[fileModel.UserID], fileModel)
		rows++
	}
	cacheMap = newCacheMap
	log.Debugf("Update %d valid records", rows)
}
