package lib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"lyn2n/event"
	"lyn2n/i18n"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
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
	timer   *time.Timer
}

func (c *Command) Exec() {
	if !c.running.TryLock() {
		log.Println("already running")
		return
	}
	defer c.running.Unlock()
	c.timer = time.NewTimer(3 * time.Minute)
	defer c.timer.Stop()
	c.cmd = c.genCmd()

	outO, err := c.cmd.StdoutPipe()
	if err != nil {
		log.Println("Error while executing command: ", err)
	}
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
		c.Stop()
	}()
	// 等待命令完成
	if err = c.cmd.Wait(); err != nil {
		log.Printf(i18n.Lang().ErrorN2NStartErr+": %v", err)
	}
	event.N2NDisConnectedEvent <- event.EmptyEvenVar
}

func (c *Command) cmdLog(outO io.ReadCloser) {
	scanner := bufio.NewScanner(outO)
	connectFlag := true
	var ip string
	connectSuccess := make(chan event.EmptySignal, 1)
	go func() {
		select {
		case <-c.timer.C:
			event.N2NConnectedErr <- event.EmptyEvenVar
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   i18n.Lang().NotifyN2NConnectErrTitle,
				Content: i18n.Lang().NotifyN2NConnectErrTimeoutContent,
			})
		case <-connectSuccess:
			c.timer.Stop()
		}
	}()
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		if strings.HasPrefix(text, "Open device") {
			ipBegin := strings.Index(text, "[ip=") + 4
			ipEnd := strings.Index(text, "][ifName")
			ip = text[ipBegin:ipEnd]
		}
		if strings.Contains(text, "Unable to set device n2n IP address") {
			event.N2NConnectedErr <- event.EmptyEvenVar
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   i18n.Lang().NotifyN2NConnectErrTitle,
				Content: i18n.Lang().NotifyN2NConnectErrTimeoutContent,
			})
			connectFlag = false
		}
		if connectFlag && strings.Contains(text, "[OK] edge <<< ================ >>> supernode") {
			event.IpChange <- ip
			event.N2NConnectedEvent <- event.EmptyEvenVar
			connectSuccess <- event.EmptyEvenVar
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   i18n.Lang().NotifyN2NConnectSuccessTitle,
				Content: i18n.Lang().NotifyN2NConnectSuccessContent + ": " + ip,
			})
		}
	}
}

func (c *Command) Stop() {
	if c.cmd == nil {
		return
	}

	if runtime.GOOS == "windows" {
		c.disConnect()
	} else {
		if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
			log.Println("Error while killing process: ", err)
		}
	}
	if len(c.StaticIp) == 0 {
		event.IpChange <- ""
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
		if err := c.cmd.Process.Signal(os.Kill); err != nil {
			log.Println("Error while killing process: ", err)
		}
	}
	if len(c.StaticIp) == 0 {
		event.IpChange <- ""
	}
}

func (c *Command) disConnect() {
	ip := "localhost"
	port := "15644"
	message := "stop"

	// 创建 UDP 地址
	addr, err := net.ResolveUDPAddr("udp", ip+":"+port)
	if err != nil {
		fmt.Println("解析地址时出错:", err)
		return
	}

	// 创建 UDP 连接
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("创建连接时出错:", err)
		return
	}
	defer conn.Close()

	// 发送消息
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("发送消息时出错:", err)
		return
	}

	fmt.Println("消息已发送:", message)
}

func (c *Command) genCmd() *exec.Cmd {
	command := exec.Command("./lib/edge", "-t", "15644", "-c", c.RoomName, "-l", c.Ip+":"+c.Port)
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
