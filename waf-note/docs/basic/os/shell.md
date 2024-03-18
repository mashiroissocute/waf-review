### shell是什么
https://wangdoc.com/bash/intro.html
shell位于操作系统kernal外层，用于操作kernal为用户提供功能。
Shell 有很多种，只要能给用户提供命令行环境的程序，都可以看作是 Shell。
历史上，主要的 Shell 有下面这些。

- Bourne Shell（sh）
- Bourne Again shell（bash）
- C Shell（csh）
- TENEX C Shell（tcsh）
- Korn shell（ksh）
- Z Shell（zsh）
- Friendly Interactive Shell（fish）
- 
通过/bin/bash打开bourne Again shell
通过/bin/sh 打开bourne shell


### mount 
https://blog.csdn.net/qq_39521554/article/details/79501714
mount挂载的作用，就是**将一个设备（通常是存储设备）挂接到一个已存在的目录上。**访问这个目录就是访问该存储设备。

### 目录相关
https://www.jianshu.com/p/bcb89e88a1be

### 调用其他shell脚本 
https://blog.csdn.net/simple_the_best/article/details/76285429

### sed 
https://www.runoob.com/linux/linux-comm-sed.html
sed 会根据脚本命令来处理文本文件中的数据，这些命令要么从命令行中输入，要么存储在一个文本文件中，此命令执行数据的顺序如下：

- 1. 每次仅读取一行内容；

- 2. 根据提供的规则命令匹配并修改数据。注意，sed 默认不会直接修改源文件数据，而是会将数据复制到缓冲区中，修改也仅限于缓冲区中的数据；
- 3. 将执行结果输出。

当一行数据匹配完成后，它会继续读取下一行数据，并重复这个过程，直到将文件中所有数据处理完毕。

http://c.biancheng.net/view/4028.html


```
replace all {TENCENTCLOUDWAF_ROOT} to ENGINE_PATH 
sed -i "s+{TENCENTCLOUDWAF_ROOT}+$WAF_ENGINE_PATH+g" $NGINX_PATH/conf/nginx.conf 

-i 直接操作源文件
s replace
g 全部替换
```
### 软连接 
https://www.cnblogs.com/sheapchen/p/015a1c8a2ebdd60d6c61e3372a5b51c0.html

### export LD_LIBRARY_PATH
https://blog.csdn.net/yyf0986/article/details/80265121

### cat >> ./sbin/nginx << EOF
https://www.jianshu.com/p/df07d8498fa5

### $@ 参数相关
https://www.cnblogs.com/fhefh/archive/2011/04/15/2017613.html

### &1 2 3 /dev/null
https://blog.csdn.net/u011630575/article/details/52151995

### autoTools
包括autoreconf、autoconf等
autotools主要就是利用各个工具来生成最后的makefile文件。其具体流程如下图:
https://geesun.github.io/posts/2015/02/autotool.html

### 动态库ldconfig LD_LIBRARY_PATH
https://www.cnblogs.com/sddai/p/10397510.html