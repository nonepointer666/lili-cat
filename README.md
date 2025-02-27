# liliCat

window 环境下的全局配置
$env:GOOS="windows"
$env:GOARCH="amd64"

$env:GOOS="linux"
编译命令
go build -o lilicat main.go

配置好密钥对后使用下面命令行上传
pscp -i "C:\Users\Administrator\.ssh\ali.ppk" "D:\code\go\lili-cat\lilicat" root@47.103.73.26:/home/admin/baique/lilicat
