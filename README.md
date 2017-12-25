# 通用日志框架


本log只是一个框架，并不处理log实际的输出  
具体的输出实现由**dirver**完成(driver是实现`log.FormatWriter`的对象)




# 获取

`go get -u github.com/afocus/log`


# 基本使用


```go

import (
    "github.com/afocus/log"
    "github.com/afocus/log/driver/console"
)


func main(){
    // 创建一个向控制台输出的log对象
    g := log.New(log.DEBUG,console.New())
    g.Info("hello,world")
    g.Error("error message")
}

```

显示结果

```
2017-12-22 13:21:43 [INFO] main:13 → user:xxx,pwd:xxx
2017-12-22 13:21:43 [ERRO] main:14 → pwd is error
```


## 输出到多个输出源

```go

import (
    "github.com/afocus/log"
    "github.com/afocus/log/driver/console"
    "github.com/afocus/log/driver/file"
)


func main(){
    // 同时输出到文件和控制台
    g := log.New(log.DEBUG,console.New(), file.New(&file.Option{...}))
    g.Info("hello,world")
    g.Error("error message")
}

```

# 高级(带有关联关系的日志)

`Ctx(id string)` 用于创建一个关联日志组

包含的方法 `Tag` 用于设置标签 支持链式调用

```go

import (
    "github.com/afocus/log"
    "github.com/afocus/log/driver/console"
)


func main(){
    // 创建一个向控制台输出的log对象
    g := log.New(log.DEBUG,console.New())

    // 创建一个日志组并关联一个id
    gx := g.Ctx("000000001").Tag("login")
    gx.Info("user:xxx,pwd:xxx")
    gx.Error("pwd is error")
    gx.Tag("newtag").Warn("dang")
    // 用完需要释放
    gx.Free()
}

```

显示结果

```
2017-12-22 13:21:43 [INFO] main:13 login-000000001 → user:xxx,pwd:xxx
2017-12-22 13:21:43 [ERRO] main:14 login-000000001 → pwd is error
```


## 输出级别

分为6级，只有当日志级别大于等于设置的级别才会输出

### 定义
Level | 说明
-----|-----
DEBUG | 指明细致的事件信息，对调试应用最有用
INFO | 指明描述信息，从粗粒度上描述了应用运行过程
WARN | 指明潜在的有害状况
ERROR | 指明错误事件，但应用可能还能继续运行
FATAL | 指明非常严重的错误事件，可能会导致应用终止执行
OFF | 最高级别，用于关闭日志

### 输出方法

* `Debug(v ...interface{})`
* `Info(v ...interface{})`
* `Warn(v ...interface{})`
* `Error(v ...interface{})`
* `Fatal(v ...interface{})`





## Driver

### 格式化输出接口

```go
// Formater 格式化日志事件到字符串
type Formater interface {
	Format(*Event) []byte
}

// FormatWriter 格式化输入
// 用于定制日志输出内容的样式
type FormatWriter interface {
	io.Writer
	Formater
}
```

实现自己的日志driver 将日志格式输出为json并通过tcp持续发送

```go

type MyDriver struct{
    c net.Conn
}

func (*MyDriver) Format(ev *log.Event)[]byte{
    b,_:=json.Marshal(ev)
    return b
}

func (d *MyDriver) Write(d []byte) int,error{
    return d.c.Write(d)
}

func main(){
    conn,err:=net.Dial("tcp","....")
    myd:=&MyDriver{c:conn}

    // 启用
    g:=log.New(log.DEBUG, myd)
    ...
}



```

### 目前已实现的driver
* [控制台](driver/console) 带颜色的控制台输出(windows暂时无色)
* [文件](driver/file) 支持文件分割
* [网络](driver/net) 目前仅实现了通过http发送日志的功能并且高度自定义







