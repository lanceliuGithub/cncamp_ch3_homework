# 项目介绍

本项目用于创建一个实验性质的HTTP服务器，仅可用于学习

# 使用说明

建议在 Linux 环境运行如下编译命令，Windows平台请先安装 Cygwin
```
make
```

运行命令
```
bin/linux/amd64/myhttpserver
```

打开方式

- 首页: http://localhost:8888
- 健康检查页: http://localhost:8888/healthz
- 缺失的页面: http://localhost:8888/no_such_page
