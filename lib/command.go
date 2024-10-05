package lib

import (
	"bufio"
	"fyne.io/fyne/v2"
	"io"
	"log"
	"lyn2n/event"
	"lyn2n/i18n"
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
		log.Println("already running")
		return
	}
	defer c.running.Unlock()
	c.cmd = c.genCmd()

	outO, err := c.cmd.StdoutPipe()
	errO, err := c.cmd.StderrPipe()
	if err != nil {
		log.Println("Error while executing command: ", err)
	}
	// 启动命令
	if err = c.cmd.Start(); err != nil {
		log.Fatalf(i18n.Lang().ErrorN2NStartErr+": %v", err)
	}
	// 使用 WaitGroup 等待异步打印
	var wg sync.WaitGroup
	// 打印标准输出
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.cmdLog(outO)
	}()

	// 打印标准错误
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(errO)
		for scanner.Scan() {
			text := scanner.Text()
			log.Println("ERROR:", text)
		}
	}()
	// 设置信号处理
	go func() {
		<-event.CloseMainWindowsEvent // 等待信号
		c.Kill()
	}()
	// 等待命令完成
	if err = c.cmd.Wait(); err != nil {
		log.Printf(i18n.Lang().ErrorN2NStartErr+": %v", err)
	}
}

func (c *Command) cmdLog(outO io.ReadCloser) {
	scanner := bufio.NewScanner(outO)
	connectFlag := true
	var ip string
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		if strings.HasPrefix(text, "Open device") {
			ipBegin := strings.Index(text, "[ip=") + 4
			ipEnd := strings.Index(text, "][ifName")
			ip = text[ipBegin:ipEnd]
		}
		if strings.Contains(text, "Unable to set device n2n IP address") {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   i18n.Lang().NotifyN2NConnectErrTitle,
				Content: i18n.Lang().NotifyN2NConnectErrContent,
			})
			connectFlag = false
		}
		if connectFlag && strings.Contains(text, "[OK] edge <<< ================ >>> supernode") {
			event.IpChange <- ip
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   i18n.Lang().NotifyN2NConnectSuccessTitle,
				Content: i18n.Lang().NotifyN2NConnectSuccessContent + ": " + ip,
			})
		}
	}
}

func (c *Command) Kill() {
	if c.cmd == nil {
		return
	}

	if runtime.GOOS == "windows" {
		if err := c.cmd.Process.Signal(os.Kill); err != nil {
			log.Println("Error while killing process: ", err)
		}
	} else {
		if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
			log.Println("Error while killing process: ", err)
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
//		log.Println("Error while loading kernel32:", err)
//		return
//	}
//	p, err := dll.FindProc("GenerateConsoleCtrlEvent")
//	if err != nil {
//		log.Println("Error while loading kernel32:", err)
//		return
//	}
//	r, _, err := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(c.cmd.Process.Pid))
//	if (err != nil && "The operation completed successfully." != err.Error()) || r == 0 {
//		log.Println("Error while loading kernel32:", err)
//	} else {
//		c.cmd = nil
//	}
//}
