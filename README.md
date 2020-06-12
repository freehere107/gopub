## gopub

Fork https://github.com/linclin/gopub 

该项目在原项目的基础上，删减大量功能，只保留一些必要的功能


## 使用框架
- [Element](http://element-cn.eleme.io/#/zh-CN)
- [Beego](https://beego.me/)
- [httprouter](https://github.com/julienschmidt/httprouter) 

## 功能特性
- Docker支持\-
- 部署简便：go二进制部署,无需安装运行环境.
- git发布支持：配置每个项目git地址,自动获取分支/tag,commit选择并自动拉取代码
- ssh执行命令/传输文件：使用golang内置ssh库高效执行命令/传输文件
- 多项目部署:支持多项目多任务并行,内置[grpool协程池](https://github.com/linclin/grpool)支持并发操作命令和传输文件
- 分批次发布：项目配置支持配置分批发布IP,自动创建多批次上线单
- 全web化操作：web配置项目,一键发布,一键快速回滚
- API支持：提供所有配置和发布操作API,便于对接其他系统  
- 部署钩子：支持部署前准备任务,代码检出后处理任务,同步后更新软链前置任务,发布完毕后收尾任务4种钩子函数脚本执行


## 部署过程

1. 创建检出目录
2. git clone repositories
3. 接收上线部署操作
4. 创建代码临时tmp_dir操作空间
5. 执行 PostDeploy 脚本
6. 执行 LastDeploy 脚本
7. clean 临时tmp_dir操作空间

