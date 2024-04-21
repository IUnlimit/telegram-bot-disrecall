package cache

import (
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
)

var EMPTY_LIST = make([]*model.FileModel, 0)

func GetTypeList(fileType model.FileType, userID int64) []*model.FileModel {
	if _, ok := cacheMap[userID]; !ok {
		return EMPTY_LIST
	}
	list := make([]*model.FileModel, 0)
	for _, fileModel := range cacheMap[userID] {
		if fileModel.FileType == fileType {
			list = append(list, fileModel)
		}
	}
	return list
}
