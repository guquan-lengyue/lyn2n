package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const edgeUrl = "https://github.com/guquan-lengyue/lyn2n/raw/refs/heads/master/lib/edge.exe"
const filePath = "./lib/edge.exe"

func checkN2NEdge() {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 文件不存在，尝试创建 lib 目录
		err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			fmt.Println("创建目录失败:", err)
			return
		}
		// 下载文件
		err = downloadFile(filePath, edgeUrl)
		if err != nil {
			fmt.Println("下载文件失败:", err)
			return
		}

		fmt.Println("文件下载成功:", filePath)
	} else {
		fmt.Println("文件已存在:", filePath)
	}
}

// downloadFile 从 URL 下载文件并保存到指定路径
func downloadFile(filepath, url string) error {
	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败: 状态码 %d", resp.StatusCode)
	}

	// 将响应体写入文件
	_, err = io.Copy(out, resp.Body)
	return err
}
