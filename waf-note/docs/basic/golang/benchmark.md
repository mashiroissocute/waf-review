基准测试（benchmark）是 go testing 库提供的，用来度量程序性能，算法优劣的利器。


benchmark 和普通的单元测试用例一样，都位于 _test.go 文件中。
函数名以 Benchmark 开头，参数是 b *testing.B。
单元测试函数名以 Test 开头，参数是 t *testing.T。


benchmark默认测试1s中函数执行的情况。
b.N 从 1 开始，如果该用例能够在 1s 内完成，b.N 的值便会增加，再次执行。b.N 的值大概以 1, 2, 3, 5, 10, 20, 30, 50, 100 这样的序列递增，越到后面，增加得越快
```golang
func BenchmarkFib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fib(30) // run fib(30) b.N times
	}
}
```


参数

- cpu : 改变使用的cpu核数
- benchtime=5s : 将默认1s修改为5s
- benchtime=50x : 设置函数执行次数为50次
- benchmem: 查看内存分配的信息


例子

```golang
package main

import (
	"math/rand"
	"testing"
	"time"
)

func generateWithCap(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func generate(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums
}

func BenchmarkGenerateWithCap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generateWithCap(1000000)
	}
}

func BenchmarkGenerate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generate(1000000)
	}
}
```

```shell
go test -benchmem -bench='.*Generate.*' . 

oos: darwin
goarch: amd64
pkg: example
BenchmarkGenerateWithCap-8  43  24335658 ns/op  8003641 B/op    1 allocs/op
BenchmarkGenerate-8         33  30403687 ns/op  45188395 B/op  40 allocs/op
```

- 43 : 1s内函数执行的次数
- 24335658 : 执行一次需要的ns数
- 8003641: 执行一次需要分配8003641B内存
- 1: 执行一次需要1次内存分配






https://www.cnblogs.com/yahuian/p/go-benchmark.html

https://geektutu.com/post/hpg-benchmark.html