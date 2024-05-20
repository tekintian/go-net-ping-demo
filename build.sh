#!/bin/sh
##################################
# 生成各个平台下的可执行程序 golang一键打包 macos, linux, windows 应用程序
# 使用方法: sh build.sh [-n appname]
# 也可忽略 -n 参数 sh build.sh  默认名称为 myapp
#
# 如: sh build.sh -n helloworld 将自动在target目录下生成以下3个可执行文件
# helloworld-darwin-amd64.bin  helloworld-linux-amd64.bin helloworld-windows-amd64.exe
#
# Author: tekintian@gmail.com
##################################

# 获取用户输入参数
while getopts ":n:" opt; do
    case $opt in
    n)
        APPNAME=$OPTARG
        ;;
    ?)
        echo "Unknown parameter"
        exit 1
        ;;
    esac
done

# -n yourappname  default app name is  myapp
APPNAME=${APPNAME:-"myapp"}
# 通用变量
export CGO_ENABLED=0 # 关闭CGO
export GOARCH=amd64  #CPU架构
# 设置darwin
export GOOS=darwin
go build -ldflags "-s -w" -o target/${APPNAME}-darwin-amd64.bin
echo "Macos可执行程序 ${APPNAME}-darwin-amd64.bin 打包成功!"
# 设置linux
export GOOS=linux
go build -ldflags "-s -w" -o target/${APPNAME}-linux-amd64.bin
echo "linux可执行程序 ${APPNAME}-linux-amd64.bin 打包成功!"
# 设置windows
export GOOS=windows
go build -ldflags "-s -w" -o target/${APPNAME}-windows-amd64.exe
echo "Windows可执行程序 ${APPNAME}-windows-amd64.exe 打包成功!"
