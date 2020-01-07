# go-idalloc

一个随笔练习的自增发号器实现

## 特性

1. 自动扩展发号器类型
2. 内存缓存自增号段，被动刷新到磁盘
3. 生产者模式利用协程通道获取结果
4. 结束进程自动更新缓存到磁盘

## Import

```golang
import (
    "github.com/chenjinya/go-idalloc"
)
```

## Demo

```golang
package main
import (
  "fmt"
  "sync"
  "time"
  "github.com/chenjinya/go-idalloc"
)

func main() {
  var ida idalloc.Pool
    idc := make(chan uint64)
    //开启一个user发号器
    ida.Run("user", idc)
    //启动自增设置为0
    ida.BootAutoIncre(0)
    //开启Debug模式
  ida.Debug(true)
  var wg sync.WaitGroup
    wg.Add(1)
    //模拟并发
  go func() {
    for {
      time.Sleep(time.Second)
      id := <-idc
      fmt.Println(id)
    }

  }()
  go func() {
    for {
      time.Sleep(time.Second)
      id := <-idc
      fmt.Println(id)
    }
  }()

  wg.Wait()
}

```

## Tips

为了防止进程重启导致id冲突，进程启动的时候，会在原来缓存的基础上自增`1000`

## License

MIT