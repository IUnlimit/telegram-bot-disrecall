package tool

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile @return path, error
func DownloadFile(url string, path string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 判断文件是否存在
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Println("检测到重复文件, 自动增加文件后缀")
		return DownloadFile(url, path+".rp")
	}

	// 获取文件所在目录的路径
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return path, nil
}
