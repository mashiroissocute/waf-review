## log
golang 自带log包
https://darjun.github.io/2020/02/07/godailylib/log

## 数据结构
``` 
// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
    mu     sync.Mutex // ensures atomic writes; protects the following fields
    prefix string     // prefix on each line to identify the logger (but see Lmsgprefix)
    flag   int        // properties
    out    io.Writer  // destination for output
    buf    []byte     // for accumulating text to write
}
```

mu : 互斥锁，防止多个GoRoutine写同个buf时，发生冲突。
buf：防止频繁的内存分配，使用了一个buf。
out：io.Writer是一个接口，有一个Write方法。 因此，out是一个实现了Write方法的对象。

写日志：
```
// 【logger.Println】等方法实际调用了log.OutPut方法。
func (l *Logger) Output(calldepth int, s string) error {
  now := time.Now() // get this early.
  var file string
  var line int
  l.mu.Lock()
  defer l.mu.Unlock()
  if l.flag&(Lshortfile|Llongfile) != 0 {
    // Release lock while getting caller info - it's expensive.
    l.mu.Unlock()
    var ok bool
    _, file, line, ok = runtime.Caller(calldepth)
    if !ok {
      file = "???"
      line = 0
    }
    l.mu.Lock()
  }
  l.buf = l.buf[:0]
  l.formatHeader(&l.buf, now, file, line)
  l.buf = append(l.buf, s...)
  if len(s) == 0 || s[len(s)-1] != '\n' {
    l.buf = append(l.buf, '\n')
  }
  _, err := l.out.Write(l.buf) 【out的write方法】
  return err
}
```


## 特性

- 可以定义log 输出的地方 ： 
例如 文件 控制台 网络 buffer等

- 可以定义log输出的前缀 （字符串）

- 可以定义log输出的格式 （时间格式 + 文件行数等）


``` 
package main

import (
    "bytes"
    "fmt"
    "io"
    "log"
    "os"
)

func main() {

    b := &bytes.Buffer{}
    f, _ := os.OpenFile("./log.txt", os.O_RDWR|os.O_CREATE, 0777)
    s := os.Stderr

    logger := log.New(io.MultiWriter(b, f, s), "logger", log.Lshortfile|log.Ltime)
    logger.Println("This is logger")

    fmt.Println(b.String())
}
```

