# syntax=docker/dockerfile:1.2

# 针对 x86_64 架构使用官方的 Golang 镜像作为基础镜像
FROM golang:1.20.8-alpine AS x86_64_builder

# 设置工作目录
WORKDIR /app

# 将项目文件复制到容器中
COPY . .

# 编译应用程序，注意替换 main.go 为您的入口文件
RUN go build -o xoj-code-sandbox main.go

# 使用 Alpine Linux 作为最终的基础镜像，针对 x86_64 架构
FROM alpine:latest AS x86_64_final

# 安装 GLIBC 和其他运行时库
RUN apk --no-cache add ca-certificates libc6-compat

# 设置工作目录
WORKDIR /app

# 复制二进制文件从构建阶段的镜像到最终的镜像
COPY --from=x86_64_builder /app/xoj-code-sandbox .

# 拷贝配置文件到容器中
COPY ./conf /app/conf

# 暴露应用程序所监听的端口
EXPOSE 8093

# 启动应用程序
CMD ["./xoj-code-sandbox"]

# 针对 ARM64 架构使用官方的 Golang 镜像作为基础镜像
FROM golang:1.20.8-alpine AS arm64v8_builder

# 设置工作目录
WORKDIR /app

# 将项目文件复制到容器中
COPY . .

# 编译应用程序，注意替换 main.go 为您的入口文件
RUN go build -o xoj-code-sandbox main.go

# 使用 Alpine Linux 作为最终的基础镜像，针对 ARM64 架构
FROM alpine:latest AS arm64v8_final

# 安装 GLIBC 和其他运行时库
RUN apk --no-cache add ca-certificates libc6-compat

# 设置工作目录
WORKDIR /app

# 复制二进制文件从构建阶段的镜像到最终的镜像
COPY --from=arm64v8_builder /app/xoj-code-sandbox .

# 拷贝配置文件到容器中
COPY ./conf /app/conf

# 暴露应用程序所监听的端口
EXPOSE 8093

# 启动应用程序
CMD ["./xoj-code-sandbox"]
