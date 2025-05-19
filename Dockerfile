# 基于 Golang 官方镜像构建
FROM golang:1.21.0-alpine3.18 AS builder

# 设置工作目录
WORKDIR /app

# 将本地应用代码复制到容器内的工作目录
COPY . .

#安装CA证书（需要请求第三方https接口）、设置代理、安装依赖、构建二进制文件
#-ldflags="-s -w":-s：省略符号表,-w：省略 DWARF 调试信息, 可进一步缩小编译后的二进制文件体积
#CGO_ENABLED=0: 强制禁用CGO，二进制文件将包含所有依赖的代码，不依赖外部动态库,允许使用 scratch 空镜像
#使用upx压缩可执行程序，能够减少程序包50%左右的体积，但会增加启动速度，需要权衡
RUN apk add --no-cache ca-certificates upx && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/main . && \
    upx --best --lzma /app/main

# 运行阶段
FROM scratch
WORKDIR /app
#复制必要文件到镜像里面
COPY . .
#复制CA证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#复制主程序
COPY --from=builder /app/main .

#设置容器暴露端口
EXPOSE 3000

CMD ["./main","serve"]