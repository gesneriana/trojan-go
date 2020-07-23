package main

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
	"encoding/base64"
	"runtime"
	"sync/atomic"

	"github.com/p4gefau1t/trojan-go/common"
	"github.com/p4gefau1t/trojan-go/log"
	"github.com/p4gefau1t/trojan-go/option"
	"github.com/p4gefau1t/trojan-go/proxy"
	_ "github.com/p4gefau1t/trojan-go/component"
)

func main() {
	// Need a main function to make CGO compile package as C shared library
}

var reentrance int64
var _proxy *proxy.Option

// Start 启动代理服务,设置可读写的配置文件的目录
//export Start
func Start(b64 *C.char) {
	var b64String = C.GoString(b64)
	decodeBytes, err := base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		log.Info(b64String)
		log.Error(err)
		return
	}

	go func() {
		if atomic.CompareAndSwapInt64(&reentrance, 0, 1) {
			log.Info("trojan-go start.")
			defer atomic.StoreInt64(&reentrance, 0)
		} else {
			log.Info("此Goroutine不可重入")
			return
		}

		var configDir = string(decodeBytes)
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

// Stop 停止客户端代理服务
//export Stop
func Stop() {
	if _proxy != nil {
		_proxy.Close()
	}
}

var platString *C.char

// GetPlatformInfo 获取系统平台
//export GetPlatformInfo
func GetPlatformInfo() *C.char {
	if platString != nil {
		return platString
	}

	platString = C.CString(runtime.GOOS + " " + runtime.GOARCH)
	return platString
}
