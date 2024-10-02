package lib

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"log"
	"lyn2n/event"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

type Command struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	RoomName string `json:"roomName"`
	RoomKey  string `json:"roomKey"`
	Encrypt  string `json:"encrypt"`
	StaticIp string `json:"staticIp"`

	cmd     *exec.Cmd
	running sync.Mutex
}

func (c *Command) Exec() {
	if !c.running.TryLock() {
		fyne.LogError("already running", nil)
		return
	}
	defer c.running.Unlock()
	c.cmd = c.genCmd()
	file, err := os.OpenFile("log.text", os.O_WRONLY|os.O_CREATE, 0644)

	outO, err := c.cmd.StdoutPipe()
	errO, err := c.cmd.StderrPipe()
	if err != nil {
		fyne.LogError("Error while executing command: ", err)
	}
	// 启动命令
	if err := c.cmd.Start(); err != nil {
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
			text := scanner.Text()
			file.WriteString(text + "\r\n")
			log.Println(text)
			if strings.HasPrefix(text, "Open device") {
				ipBegin := strings.Index(text, "[ip=") + 4
				ipEnd := strings.Index(text, "][ifName")
				event.IpChange <- text[ipBegin:ipEnd]
			}
			if strings.Contains(text, "> nul]") {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "n2n连接失败",
					Content: "无法给虚拟网卡设置ip",
				})
			}
		}
	}()

	// 打印标准错误
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(errO)
		for scanner.Scan() {
			text := scanner.Text()
			file.WriteString(text + "\r\n")
			fmt.Println("ERROR:", text)
		}
	}()
	// 设置信号处理
	go func() {
		<-event.CloseMainWindowsEvent // 等待信号
		c.Kill()
	}()
	// 等待命令完成
	if err := c.cmd.Wait(); err != nil {
		log.Printf("命令执行失败: %v", err)
	}

}

func (c *Command) Kill() {
	if c.cmd == nil {
		return
	}

	if runtime.GOOS == "windows" {
		if err := c.cmd.Process.Signal(os.Kill); err != nil {
			fyne.LogError("Error while killing process: ", err)
		}
	} else {
		if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
			fyne.LogError("Error while killing process: ", err)
		}
	}
	if len(c.StaticIp) == 0 {
		event.IpChange <- ""
	}
}

func (c *Command) genCmd() *exec.Cmd {
	command := exec.Command("./lib/edge", "-c", c.RoomName, "-l", c.Ip+":"+c.Port)
	if len(c.RoomKey) > 0 {
		command.Args = append(command.Args, "-k", c.RoomKey, c.encryptCmd())
	}
	if len(c.StaticIp) > 0 {
		command.Args = append(command.Args, "-a", c.StaticIp)
	}
	command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
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

//
//func (c *Command) windowKill() {
//	dll, err := syscall.LoadDLL("kernel32.dll")
//	if err != nil {
//		fyne.LogError("Error while loading kernel32:", err)
//		return
//	}
//	p, err := dll.FindProc("GenerateConsoleCtrlEvent")
//	if err != nil {
//		fyne.LogError("Error while loading kernel32:", err)
//		return
//	}
//	r, _, err := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(c.cmd.Process.Pid))
//	if (err != nil && "The operation completed successfully." != err.Error()) || r == 0 {
//		fyne.LogError("Error while loading kernel32:", err)
//	} else {
//		c.cmd = nil
//	}
//}
