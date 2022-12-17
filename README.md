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

# 制作容器镜像

生成容器镜像
```
make release
```

生成容器镜像并推送到 Docker Hub 公开仓库
```
make push
```

如果推送时报错 `denied: requested access to the resource is denied` ，请先登录 docker.com
```
docker login
```

# 启动应用（请提前安装好Docker）
```
docker run -d -e VERSION=1.0 -p 80:8888 lanceliu2022/myhttpserver:1.0
```

打开方式

- 首页: http://localhost:8888
- 健康检查页: http://localhost:8888/healthz
- 缺失的页面: http://localhost:8888/no_such_page
