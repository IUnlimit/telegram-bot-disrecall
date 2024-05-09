package tool

import (
	"io"

	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
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
		log.Debug("Duplicate files are detected and the file suffix is automatically added")
		return DownloadFile(url, path+".rp")
	}

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

func CopyFile(src string, target string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	dir := filepath.Dir(target)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	destFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}
