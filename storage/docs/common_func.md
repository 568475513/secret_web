## 公共方法
公共方法位置：pkg/util
```markdown
//utils.go(公共函数)
    - rootPath() //返回项目根目录
    - Struct2Map(obj interface{}) map[string]interface{} //结构体转化为字典
    - StructJsonMap(obj interface{}, v *map[string]interface{}) error //结构体Json转化为字典
    - GetRuntimeDir() string //获取runtime目录路径
    - IsQyApp(versionType string) bool //是否是app
```