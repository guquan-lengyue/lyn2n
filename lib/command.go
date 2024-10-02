package lib

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"log"
	"os/exec"
	"sync"
)

type Command struct {
	Ip       string
	Port     string
	RoomName string
	RoomKey  string
	Encrypt  string
}

func Exec(cmd *Command) {
	command := cmd.genCmd()
	outO, err := command.StdoutPipe()
	errO, err := command.StderrPipe()
	if err != nil {
		fyne.LogError("Error while executing command: ", err)
	}
	// 启动命令
	if err := command.Start(); err != nil {
		log.Fatalf("命令启动失败: %v", err)
	}
	// 使用 WaitGroup 等待异步打印
	var wg sync.WaitGroup
	// 打印标准输出
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(outO)
		for scanner.Scan() {
			fmt.Println("OUTPUT:", scanner.Text())
		}
	}()

	// 打印标准错误
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(errO)
		for scanner.Scan() {
			fmt.Println("ERROR:", scanner.Text())
		}
	}()
	// 等待命令完成
	if err := command.Wait(); err != nil {
		log.Fatalf("命令执行失败: %v", err)
	}

}

func (c *Command) genCmd() *exec.Cmd {
	command := exec.Command("./lib/edge", "-c", c.RoomName, "-l", c.Ip+":"+c.Port)
	if len(c.RoomKey) > 0 {
		command.Args = append(command.Args, "-k", c.RoomKey, c.encryptCmd())
	}
	return command
}

func (c *Command) encryptCmd() string {
	switch c.Encrypt {
	case "Twofish":
		return "-A2"
	case "AES":
		return "-A3"
	case "ChaCha20":
		return "-A4"
	case "Speck-CTR":
		return "-A5"
	}
	return ""
}
