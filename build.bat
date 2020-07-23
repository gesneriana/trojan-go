echo 请确保你已经安装go 1.12+ 并且成功安装gcc (Windows中是MingW)
go build -buildmode=c-shared -tags "full" -ldflags="-s -w" -o trojan_go_64.dll main.go