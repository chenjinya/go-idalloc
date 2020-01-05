# go-idalloc

一个随笔练习的自增发号器实现

## 特性

1. 自动扩展发号器类型
2. 内存缓存自增号段，被动刷新到磁盘

## USAGE

```
import (
    "github.com/chenjinya/go-idalloc"
)
```

## DEMO

```golang

package main

import (
    "fmt"
    "github.com/chenjinya/go-idalloc"
)

func main() {

    ida := idalloc.Idalloc{
        Type: "user"}

    //设置为Debug模式，展示日志
    ida.Debug(false)

    //将自动自动增加步长设为0
    ida.BootAutoIncre(0)

    //获得ID
    id, err := ida.Gen()

    if nil != err {
        fmt.Println(err)
    }

    fmt.Println(id)
}

```

## Tips

为了防止进程重启导致id冲突，进程启动的时候，会在原来缓存的基础上自增`1000`

## License

MIT