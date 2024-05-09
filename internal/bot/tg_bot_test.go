package bot

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	fileUID := "=qwer="
	// 测试替换文件名为指定的fileUID
	fmt.Println(generateFileName("photos/file_0", fileUID))                                         // 输出：/path/png/file123
	fmt.Println(generateFileName("/var/lib/telegram-bot-api/<token>/photos/file_464", fileUID))     // 输出：/path/png/file123.jpg
	fmt.Println(generateFileName("/var/lib/telegram-bot-api/<token>/photos/file_464.avi", fileUID)) // 输出：/path/png/file123.mp4
}
