package trojangolib

import (
	"runtime"
	"sync/atomic"

	"github.com/p4gefau1t/trojan-go/common"
	"github.com/p4gefau1t/trojan-go/log"
	"github.com/p4gefau1t/trojan-go/option"
	"github.com/p4gefau1t/trojan-go/proxy"

		_ "github.com/p4gefau1t/trojan-go/component"
)

var reentrance int64
var _proxy *proxy.Option

type SetGeoFileFunc func(dir string)

var SetGeoFile SetGeoFileFunc

// Start the proxy service and set the directory of readable and writable configuration files
//export Start
func Start(configRootDir string) {

	go func() {
		if atomic.CompareAndSwapInt64(&reentrance, 0, 1) {
			log.Info("trojan-go start.")
			defer atomic.StoreInt64(&reentrance, 0)
		} else {
			log.Info("此Goroutine不可重入")
			return
		}

		var configFile = configRootDir + "/config.json"
		// jni 调用golib只会执行一次初始化, 所以再次启动需要手动初始化
		if option.HandlersCount() == 0 {
			var o = &proxy.Option{}
			o.SetConfigJsonPath(configFile)
			option.RegisterHandler(o)
		}

		common.ProgramDir = configRootDir

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

// Stop the client proxy service
//export Stop
func Stop() {
	if _proxy != nil {
		_proxy.Close()
	}
}

// GetPlatformInfo 获取系统平台
//export GetPlatformInfo
func GetPlatformInfo() string {
	return runtime.GOOS + " " + runtime.GOARCH
}
