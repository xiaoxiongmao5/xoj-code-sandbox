# xoj-code-sandbox

## 部署步骤

📢 注意：该项目目前不支持 Docker 容器部署，因为 Docker 容器内的服务，不能直接访问到宿主机的 Docker Daemon。
1. 打包：
  ```bash
  go build
  ```
2. 下载镜像：
  ```bash
  docker pull golang:1.20.8-alpine
  docker pull alpine:latest
  ```
3. 后台运行：
  ```bash
  nohup ./xoj-code-sandbox &
  ```
