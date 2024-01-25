介绍:https://www.liwenzhou.com/posts/Go/singleflight/

原理:https://www.lixueduan.com/posts/go/singleflight/

singleflight属于Go的准标准库，它提供了重复函数调用抑制机制，使用它可以避免同时进行相同的函数调用。第一个调用未完成时后续的重复调用会等待，当第一个调用完成时则会与它们分享结果，这样以来虽然只执行了一次函数调用但是所有调用都拿到了最终的调用结果。

singleflight内部使用 waitGroup 来让同一个 key 的除了第一个请求的后续所有请求都阻塞。直到第一个请求执行 fn 返回后，其他请求才会返回。


使用注意点:

- singlefilght的第一个请求是阻塞的，如果第一个请求一直拿不到数据，那么后续的请求也会被全部阻塞.解决方法是，使用 DoChan 结合 ctx + select 做超时控制。

- singleflight第一个请求失败了，那么后续所有等待的请求都会返回同一个 error。
实际上可以根据下游能支撑的rps定时forget一下key，让更多的请求能有机会走到后续逻辑。