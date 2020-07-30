echo 请确保你已经安装go 1.12+ 并且成功安装android ndk 和 gomobile等工具
gomobile bind -target=android -o ./golibs.aar -ldflags="-s -w" -tags="full"