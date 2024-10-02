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
	"sync"
	"syscall"
)

type Command struct {
	Ip       string
	Port     string
	RoomName string
	RoomKey  string
	Encrypt  string

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
	// 设置信号处理
	go func() {
		<-event.CloseMainWindowsEvent // 等待信号
		fmt.Println("收到停止信号，正在停止命令...")
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
		dll, err := syscall.LoadDLL("kernel32.dll")
		if err != nil {
			fyne.LogError("Error while loading kernel32:", err)
			return
		}
		p, err := dll.FindProc("GenerateConsoleCtrlEvent")
		if err != nil {
			fyne.LogError("Error while loading kernel32:", err)
			return
		}
		r, _, err := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(c.cmd.Process.Pid))
		if (err != nil && "The operation completed successfully." != err.Error()) || r == 0 {
			fyne.LogError("Error while loading kernel32:", err)
		} else {
			c.cmd = nil
		}
	} else {
		if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
			fyne.LogError("Error while killing process: ", err)
		}
	}
}

func (c *Command) genCmd() *exec.Cmd {
	command := exec.Command("./lib/edge", "-c", c.RoomName, "-l", c.Ip+":"+c.Port)
	if len(c.RoomKey) > 0 {
		command.Args = append(command.Args, "-k", c.RoomKey, c.encryptCmd())
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
