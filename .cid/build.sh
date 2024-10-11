#!/bin/bash
app="ats"
[ -d output ] || mkdir output
# 比对项目文件中引入的依赖与go.mod进行比对,清理不需要的依赖,并且更新go.sum文件。
go mod tidy
# 将项目的所有依赖导出至vendor目录
go mod vendor
# 构建
LD_PATH="ats/src/config"
APIGW_VERSION="0.0.1"
GO_VERSION=$(go version |awk '{print $3}')
REVISION="master"
LD_FLAGS="-X $LD_PATH.Version=${APIGW_VERSION} -X $LD_PATH.GoVersion=${GO_VERSION} -X $LD_PATH.GitCommit=${REVISION}"
go build  -trimpath -ldflags="-s -w $LD_FLAGS" -gcflags='all=-l -N' -o output/$app src/main/main.go
