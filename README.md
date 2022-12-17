# 项目介绍

本项目用于创建一个实验性质的HTTP服务器，仅可用于学习

# 编译二进制可执行文件

建议在 Linux 环境运行如下编译命令，Windows平台请先安装 Cygwin
```
make
```
或
```
make build
```

运行命令
```
bin/linux/amd64/myhttpserver
```

打开方式

- 首页: http://localhost:8888
- 健康检查页: http://localhost:8888/healthz
- 缺失的页面: http://localhost:8888/no_such_page

# 制作容器镜像

生成容器镜像
```
make release
```

生成容器镜像并推送到 Docker Hub 公开仓库
```
make push
```
