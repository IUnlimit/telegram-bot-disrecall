package cache

import (
	"github.com/IUnlimit/telegram-bot-disrecall/internal/model"
)

var EmptyList = make([]*model.FileModel, 0)

// GetTypeList 获取指定文件类型的数据 list
func GetTypeList(fileType model.FileType, userID int64) []*model.FileModel {
	if _, ok := cacheMap[userID]; !ok {
		return EmptyList
	}

	// MediaGroup 全返回
	if fileType == model.MediaGroup {
		return cacheMap[userID]
	}

	list := make([]*model.FileModel, 0)
	for _, fileModel := range cacheMap[userID] {
		if fileModel.FileType == fileType {
			list = append(list, fileModel)
		}
	}
	return list
}

// GetUserStatic 获取用户储存的静态数据
// map - fileType: MB
func GetUserStatic(userID int64) map[model.FileType]*model.Static {
	static := make(map[model.FileType]*model.Static)
	types := []model.FileType{model.Text, model.Voice, model.Photo, model.Video, model.Document}
	for _, fileType := range types {
		s := model.Static{}
		static[fileType] = &s
	}
	if _, ok := cacheMap[userID]; !ok {
		return static
	}

	for _, fileModel := range cacheMap[userID] {
		s := static[fileModel.FileType]
		s.Rows += 1
		s.MB += float64(fileModel.FileSize) / 1024 / 1024
	}
	return static
}
