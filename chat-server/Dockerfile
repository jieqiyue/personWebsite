# 第一步：构建阶段，使用 Go 官方镜像编译
FROM golang:1.19 as builder

WORKDIR /app

# 把代码复制进来
COPY . .

# 拉取依赖并编译
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o app

# 第二步：运行阶段，使用一个超小的基础镜像
FROM alpine:latest

WORKDIR /

# 只把编译好的app程序拷贝过来
COPY --from=builder /app/app .

# 映射需要监听的端口，比如你的go项目是8080
EXPOSE 8080

# 启动
ENTRYPOINT ["./app"]
