## 背景
TencentEngine中使用，ngx.re.find, ngx.re.match, ngx.re.gmatch来进行正则匹配。
## ngx.re.match
captures, err = ngx.re.match(subject, regex, options?, ctx?, res_table?)  // ?为可选

**含义：** Matches the  `subject`  string using the Perl compatible regular expression  `regex`  with the optional  `options` .

**captures：** 未匹配上时，captures为nil；匹配上时，captures为lua table。 captures[0]返回`第一个`匹配上的子串。captures[1]返回正则表达式中第一个**括号**匹配上的子串，captures[2]返回正则表达式中第二个**括号**匹配上的子串... 如果括号内正则没有匹配上，则返回false。
```
 local m, err = ngx.re.match("hello, 1234", "([0-9])([0-9]+)")
 -- m[0] == "1234"
 -- m[1] == "1"
 -- m[2] == "234"

 local m, err = ngx.re.match("hello, world", "(world)|(hello)|(?<named>howdy)")
 -- m[0] == "hello"
 -- m[1] == false
 -- m[2] == "hello"
 -- m[3] == false
 -- m["named"] == false //支持为括号内正则取名字，通过key访问
```
**err：** 当正则匹配过程中，出现错误，通过err string 描述。错误例如，超出PCRE栈大小限制，或者错误的正则表达式语法。
**options: ** 正则匹配模式选项
jo 一起使用，利于性能调优。 
The  `o`  option is useful for performance tuning, because the regex pattern in question will only be compiled once, cached in the worker-process level, and shared among all requests in the current Nginx worker process. 一次编译，存放在worker cache，所有的请求通用。
```
j             enable PCRE JIT compilation, this requires PCRE 8.21+ which
              must be built with the --enable-jit option. for optimum performance,
              this option should always be used together with the 'o' option.
              first introduced in ngx_lua v0.3.1rc30.

o             compile-once mode (similar to Perl's /o modifier),
              to enable the worker-process-level compiled-regex cache

s             single-line mode (similar to Perl's /s modifier)

u             UTF-8 mode. this requires PCRE to be built with
              the --enable-utf8 option or else a Lua exception will be thrown.
                                         ......
```
**ctx：** lua table，指定从字符串的什么位置开始匹配。该值在匹配过程中不不断变化
```
 local ctx = { pos = 2 }
 local m, err = ngx.re.match("1234, hello", "[0-9]+", "", ctx)
      -- m[0] = "234"
      -- ctx.pos == 5
```

## ngx.re.find
from, to, err = ngx.re.find(subject, regex, options?, ctx?, nth?)
**含义：** returns the beginning index ( `from` ) and end index ( `to` ) of the matched substring. The returned indexes are 1-based and can be fed directly into the [string.sub](https://www.lua.org/manual/5.1/manual.html#pdf-string.sub) API function to obtain the matched substring.
**from to：** index，默认参数下，会返回`第一个`匹配上子串的index。
**options ctx：** 和match语意一致
**nth：** 返回正则表达式中第n个**括号**匹配上的子串，类似capture的index。

## ngx.re.gmatch
iterator, err = ngx.re.gmatch(subject, regex, options?)
**含义：** returns a Lua iterator instead, so as to let the user programmer iterate all the matches over the  `<subject>`  string argument with the PCRE  `regex` .通过迭代器返回所有匹配上的子串。不再是第一个子串了。
**示例：**
```
 local it, err = ngx.re.gmatch("hello, world!", "([a-z]+)", "i")
 if not it then
     ngx.log(ngx.ERR, "error: ", err)
     return
 end

 while true do
     local m, err = it() //这里会从it中取出数据
     if err then
         ngx.log(ngx.ERR, "error: ", err)
         return
     end

     if not m then //it取完数据后，m为nil，退出循环
         -- no match found (any more)
         break
     end

     -- found a match
     ngx.say(m[0]) 
     ngx.say(m[1])
 end
```

## ngx.re.sub
newstr, n, err = ngx.re.sub(subject, regex, replace, options?)
使用relace替换`第一个`匹配上的子串
replace中可以使用$来使用正则匹配子串，例如`$0`  referring to the whole substring matched by the pattern and  `$1`  referring to the first parenthesized(括号) capturing substring.
```
 local newstr, n, err = ngx.re.sub("hello, 1234", "([0-9])[0-9]", "[$0][$1]")
 if not newstr then
     ngx.log(ngx.ERR, "error: ", err)
     return
 end

 -- newstr == "hello, [12][1]34"
 -- n == 1
```
replace还可以是function
```
 local func = function (m)
     return "[" .. m[0] .. "][" .. m[1] .. "]"
 end

 local newstr, n, err = ngx.re.sub("hello, 1234", "( [0-9] ) [0-9]", func, "x")
 -- newstr == "hello, [12][1]34"
 -- n == 1
```
## ngx.re.gsub
使用relace替换`所有`匹配上的子串
```
 local newstr, n, err = ngx.re.gsub("hello, world", "([a-z])[a-z]+", "[$0,$1]", "i")
 if not newstr then
     ngx.log(ngx.ERR, "error: ", err)
     return
 end

 -- newstr == "[hello,h], [world,w]"
 -- n == 2
```


## 参考
https://github.com/openresty/lua-nginx-module