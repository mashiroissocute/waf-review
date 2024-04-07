## Golang版本管理
类似Conda用来管理Python版本。Golang官方推出了golang的版本管理。
Golang不同版本的区别在哪里？ golang的编译器，runtime等代码都会发生改变。

## 具体操作
https://polarisxu.studygolang.com/posts/go/managing-multiple-go-versions/

go env 找到GOPATH
进入GOPATH，执行： `go get golang.org/dl/go<version>` <version>是可变的版本号
进入GOPATH/bin 执行 `./go<version> download` 下载
进入项目执行 GOPATH/bin/go<version> mod init xxx 
编译： GOPATH/bin/go<version> build ...
