# March

使用 go+vue 开发的一个 简易 linux 系统监控

## 已实现

- CPU 使用率监控
- 内存使用率监控
- SWAP 使用率监控
- 进程监控
- 磁盘 IO 监控
- 网络负载

## 使用方法

```sh
yarn run build  # 构建前端
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' .  # 构建后端
```

## 页面截图

![截图](./screenshots.jpeg)
