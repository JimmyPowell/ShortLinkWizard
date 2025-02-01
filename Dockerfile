# 第一阶段：构建 Go 后端
FROM golang:1.20 AS go-builder

WORKDIR /app

# 复制 Go 模块文件
COPY go.mod go.sum ./

# 下载依赖（修正命令）
RUN go mod tidy

# 复制所有项目文件
COPY . .

# 构建可执行文件（修正环境变量和参数）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

# 第二阶段：构建 React 前端
FROM node:18 AS react-builder

WORKDIR /app

# 复制 package 文件（修正路径）
COPY front_end/package*.json .

# 安装依赖
RUN npm install

# 复制前端项目文件（修正路径）
COPY front_end .

# 构建 React 项目
RUN npm run build

# 第三阶段：创建轻量级镜像
FROM alpine:latest

# 安装依赖（修正命令顺序）
RUN apk add --no-cache ca-certificates

WORKDIR /root

# 复制 Go 可执行文件
COPY --from=go-builder /app/main .

# 复制 React 静态文件（假设构建路径正确）
COPY --from=react-builder /app/build ./static

EXPOSE 8081

# 运行应用（修正 CMD 格式）
CMD ["./main"]
