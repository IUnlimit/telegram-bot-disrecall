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
	timer.AddFunc("@every 5m", func() {
		// TODO 没有插入记录的不刷新
		updateLocalCache()
	})
	timer.Start()
}

func updateLocalCache() {
	cacheMap = make(map[int64][]*model.FileModel, 0)
	fileModels, err := db.QueryAll()
	if err != nil {
		log.Fatalf("Initial cache data failed: %v", err)
	}

	var rows int
	for _, fileModel := range fileModels {
		if _, ok := cacheMap[fileModel.ForwardUserID]; !ok {
			cacheMap[fileModel.ForwardUserID] = make([]*model.FileModel, 0)
		}
		if !fileModel.IsValid() {
			continue
		}
		cacheMap[fileModel.ForwardUserID] = append(cacheMap[fileModel.ForwardUserID], fileModel)
		rows++
	}
	log.Debugf("Update %d valid records", rows)
}
