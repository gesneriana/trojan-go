package trojangolib

import (
	"net"
	"runtime"
	"sync/atomic"

	"github.com/p4gefau1t/trojan-go/common"
	_ "github.com/p4gefau1t/trojan-go/component"
	"github.com/p4gefau1t/trojan-go/log"
	"github.com/p4gefau1t/trojan-go/option"
	"github.com/p4gefau1t/trojan-go/proxy"
	"github.com/p4gefau1t/trojan-go/tun2socks"
)

func main() {
	// Need a main function to make CGO compile package as C shared library
}

var reentrance int64
var _proxy *proxy.Option

// Start trojan proxy, set config dir
//export Start
func Start(configDir string) {

	go func() {
		if atomic.CompareAndSwapInt64(&reentrance, 0, 1) {
			log.Info("trojan-go start.")
			defer atomic.StoreInt64(&reentrance, 0)
		} else {
			log.Info("此Goroutine不可重入")
			return
		}

		common.ProgramDir = configDir
		var configFile = configDir + "/config.json"
		// jni 调用golib只会执行一次初始化, 所以再次启动需要手动初始化
		if option.HandlersCount() == 0 {
			var o = &proxy.Option{}
			o.SetConfigJsonPath(configFile)
			option.RegisterHandler(o)
		}
		// 不能使用defer释放锁，因为下面代码会阻塞，造成死锁

		for {
			h, err := option.PopOptionHandler()
			if err != nil {
				// 为了android平台加入了手动关闭的机制, 防止golib导致app闪退
				log.Warn("invalid options", err)
				break
			}

			v, ok := h.(*proxy.Option)
			if ok {
				_proxy = v
				_proxy.SetConfigJsonPath(configFile)
				log.Info("SetConfigJsonPath: " + configFile)
			}

			// log.Info(reflect.TypeOf(h).String())

			err = h.Handle()
			if err == nil {
				break
			}
		}
	}()

}

// Stop trojan proxy
//export Stop
func Stop() {
	if _proxy != nil {
		_proxy.Close()
	}
}

// Tun2socksStartOptions is tun2socks.Tun2socksStartOptions
type Tun2socksStartOptions struct {
	TunFd        int
	Socks5Server string
	FakeIPRange  string
	MTU          int
	EnableIPv6   bool
}

// StartTun sets up lwIP stack, starts a Tun2socks instance
//export StartTun
func StartTun(opt *Tun2socksStartOptions) int {
	var tunopt = &tun2socks.Tun2socksStartOptions{}
	tunopt.EnableIPv6 = opt.EnableIPv6
	tunopt.FakeIPRange = opt.FakeIPRange
	tunopt.MTU = opt.MTU
	tunopt.Socks5Server = opt.Socks5Server
	tunopt.TunFd = opt.TunFd
	return tun2socks.Start(tunopt)
}

// StopTun stop tun
//export StopTun
func StopTun() {
	tun2socks.Stop()
}

// GetFreePort asks the kernel for a free open port that is ready to use.
//export GetFreePort
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// GetPlatformInfo get system pltform
//export GetPlatformInfo
func GetPlatformInfo() string {
	return runtime.GOOS + " " + runtime.GOARCH
}
