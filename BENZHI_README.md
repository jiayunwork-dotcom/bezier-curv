# bezier-curv：基于 Go 实现的三次 Bézier 曲线弧长、曲率与法向等距偏移核算 Web 服务

输入四个控制点（JSON），通过 HTTP API 或 CLI 子命令，用 Bernstein 基与数值积分输出弧长（自适应 Simpson 收敛值）、曲率 κ(t) 与法向偏移折线；非法输入与尖点会明确报错。服务默认监听 :8080，提供 /health、/api/arclength、/api/curvature、/api/offset、/api/validate、/api/check、/api/sample 等路由。

## 构建 / 运行 / 测试

```text
go build ./...
go run . sample example/cubic-bow.json
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
