# AbsGo 开发规范（v1）
> 修订人：Dexterli

abs-go是不能在windows环境跑的，如果要跑，有以下方法：
- 修改main.go去掉github.com/fvbock/endless 包的使用。（里面有说明，但是不要推修改的代码到开发环境）
- 打开Microsoft Store 搜索Terminal安装，然后在搜索Ubuntu安装，即可在windows打开Ubuntu运行。

## 命名规范
- 全局常量命名请全部大写，语义下划线隔开，参考  `util.TIME_LAYOUT`

- 文件包内常量和变量遵循GO语言设计模式，私有开头小写，公有开头大写，注意是**驼峰法命名**

- 文件和文件夹命名使用**蛇形命名**，语义下划线隔开

- 普通变量也请遵循**驼峰法命名**，开头小写

- 类方法和方法遵循GO语言设计模式，采用**驼峰法命名**

- 各类命名请考虑到**语义**，比如直播详情方法：GetAliveInfo()
## 注释
- 所有方法都要加上注释，需要有一个大概的解释！！注释为行注释即可
- util里面的方法需要在common_func.md补充解释说明
## 包管理
请用括号模式
```golang
import (
    "fmt"
    ...
)
```
多种类型的包请按照
```golang
import (
    // 系统标准库包
	"encoding/json"
	"fmt"
	"time"

    // 引用第三方包
	"github.com/gomodule/redigo/redis"

    // 内部包
	"abs/models/alive"
	"abs/pkg/cache/alive_static"
	"abs/pkg/cache/redis_alive"
	"abs/pkg/logging"
)
```
这样层次引入，注意：**空行隔开**
## 开发流程规定