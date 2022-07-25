FROM golang:alpine
  
# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

RUN go build -o app .

EXPOSE 9205

CMD ["./app"]