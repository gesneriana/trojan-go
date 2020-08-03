# Trojan-Go


## **下面的说明为简单介绍，完整配置教程和配置介绍参见[Trojan-Go文档](https://p4gefau1t.github.io/trojan-go)。**

此分支为移植到移动端的代码, 请参考主分支的readme文档

## 构建

确保你的Go版本 >= 1.14，推荐使用snap安装Go保持与上游同步。
请按照Android SDK, NDK, gomobile, JAVA JDK

下面的命令使用gomobile进行编译

```shell
git clone https://github.com/p4gefau1t/trojan-go.git
cd trojan-go
运行build.bat文件或者执行命令: gomobile bind -target=android -o ./golibs.aar -ldflags="-s -w" -tags="full"
运行成功后会得到Android的aar包, android项目可以直接引用, flutter项目需要以flutter插件形式使用

## 致谢

[trojan-go](https://github.com/p4gefau1t/trojan-go)

[trojan](https://github.com/trojan-gfw/trojan)

[v2ray](https://github.com/v2ray/)

[smux](https://github.com/xtaci/smux)

[go-tproxy](https://github.com/LiamHaworth/go-tproxy)

[utls](https://github.com/refraction-networking/utls)

## Stargazers over time

[![Stargazers over time](https://starchart.cc/p4gefau1t/trojan-go.svg)](https://starchart.cc/p4gefau1t/trojan-go)
