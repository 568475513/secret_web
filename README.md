## 目录结构
```
├── cmd                  项目服务入口（internal可调用入口，静态文件或路由配置类可放在此层）
│   ├── job              项目异步脚本任务服务
│   └── server           web处理服务
│       └── routers      路由层
│           └── groups   路由组文件，路由需要分各种文件，不写在一个路由文件
├── internal             服务私有化局域包（内部互不调用）
│   ├── job              项目异步脚本任务服务
│   │   ├── repository   业务逻辑仓库层（不可互相调用）
│   │   └── tasks        异步任务处理入口（类似于控制器）
│   └── server           web处理服务
│       ├── api/v2       控制器（不可使用Model层、Redis层、Server层等数据层，数据和大部分业务逻辑请在仓库层实现）
│       ├── middleware   中间件
│       ├── repository   业务逻辑仓库层（不可互相调用）
│       └── rules        请求参数鉴权层
├── pkg                  功能工具包库（不可使用外部的自建服务包，不可写入业务代码！！！）
│   ├── app              请求返回处理包库，比如格式返回和请求参数校验等
│   ├── cache            缓存库，比如redis连接池的初始化
│   ├── conf             配置库，配置初始化
│   ├── enums            标准常量库，即类型、以及返回code的常量设置包
│   ├── file             文件处理库
│   ├── logging          日志处理库
│   ├── provider         服务重写库，比如对SQLNULLString的业务重写
│   └── util             工具库，类似php的helpers
├── runtime              运行输出文件夹
│   └── logs             程序里面添加的输出日志
├── service              服务请求数据层，封装对各种第三方服务的请求（比如权益黑名单等）处理库（只可引入PKG工具包）
├── models               数据库数据模型层，只负责取数据库数据层（只可引入PKG工具包）
├── vendor               加载的库包，类似php的
```