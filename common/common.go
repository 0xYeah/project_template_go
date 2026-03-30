package common

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/0xYeah/project_template_go/config"
	"github.com/george012/gtbox"
	"github.com/george012/gtbox/gtbox_coding"
	"github.com/george012/gtbox/gtbox_log"
	"github.com/george012/gtbox/gtbox_net"
	"github.com/george012/gtbox/gtbox_time"
)

var (
	globalTransport *http.Transport
	initOnce        sync.Once
)

func GetGlobalHttpClient(requestTimeout time.Duration) *http.Client {
	initOnce.Do(func() {
		globalTransport = &http.Transport{
			MaxIdleConns:          500,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   requestTimeout,
			ResponseHeaderTimeout: requestTimeout,
			ExpectContinueTimeout: 1 * time.Second,
		}
	})
	return &http.Client{
		Transport: globalTransport,
		Timeout:   requestTimeout,
	}
}

func LoadSigHandle(cleanAction func(), testMethods []func()) {
	if config.CurrentApp != nil && config.CurrentApp.CurrentRunMode == gtbox.RunModeDebug {
		testMethod(testMethods)
	} else if config.CurrentApp == nil {
		testMethod(testMethods)
	}
	// 创建一个信号通道
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)

	// 阻断主进程等待signal
	asig := <-chSig
	if cleanAction != nil {
		cleanAction()
	}
	gtbox_log.LogInfof("接收到 [%s] 信号，程序即将退出! ", asig)
	willExitHandle()
}

// willExitHandle 异常退出处理
func willExitHandle() {
	gtbox_log.LogInfof("[程序关闭]---[处理缓存数据] ")

	// 退出
	ExitApp()
}

func testMethod(testMethods []func()) {
	line_No := gtbox_coding.GetProjectCodeLines()
	gtbox_log.LogDebugf("项目有效代码总行数: %v", line_No)
	gtbox_log.LogDebugf("当前公网IP: %v", gtbox_net.GTGetPublicIPV4())
	for _, method := range testMethods {
		go method()
		gtbox_log.LogDebugf("开始执行测试方法: %v", method)

	}
}

func ExitApp() {
	// 发送 os.Interrupt 信号以触发正常退出
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(os.Interrupt)
}

func PanicHandler(context ...string) {
	err := recover()
	if err != nil {
		timestamp := gtbox_time.NowUTC().Format("2006-01-02T15:04:05.000") + " UTC"
		gtbox_log.LogErrorf("====================== Panic Error ======================")
		if len(context) > 0 {
			gtbox_log.LogErrorf("[%s] 上下文: %s", timestamp, strings.Join(context, ", "))
		}
		gtbox_log.LogErrorf("程序遇到严重错误")
		gtbox_log.LogErrorf("错误详细信息: %+v", err)

		// 获取堆栈信息
		stack := debug.Stack()
		stackLines := strings.Split(string(stack), "\n")

		// 提取触发 panic 的方法名、文件和行号
		var methodName, fileName string
		var lineNum int
		for i, line := range stackLines {
			if strings.Contains(line, ".go:") && !strings.Contains(line, "runtime/") {
				// 假设第一行非 runtime 的堆栈帧是触发点
				if i > 0 && i-1 < len(stackLines) {
					prevLine := stackLines[i-1]
					if idx := strings.LastIndex(prevLine, "("); idx > 0 {
						methodName = strings.TrimSpace(prevLine[:idx])
						if lastDot := strings.LastIndex(methodName, "."); lastDot > 0 {
							methodName = methodName[lastDot+1:]
						}
					}
				}
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					fileName = strings.TrimSpace(parts[0])
					if colonIdx := strings.Index(parts[1], " "); colonIdx > 0 {
						lineNumStr := parts[1][:colonIdx]
						fmt.Sscanf(lineNumStr, "%d", &lineNum)
					}
				}
				break
			}
		}

		gtbox_log.LogErrorf("触发 Panic 的方法: %s", methodName)
		gtbox_log.LogErrorf("文件: %s, 行号: %d", fileName, lineNum)
		gtbox_log.LogErrorf("完整堆栈跟踪:\n%s", string(stack))
		gtbox_log.LogErrorf("=========================================================")
	}
}
