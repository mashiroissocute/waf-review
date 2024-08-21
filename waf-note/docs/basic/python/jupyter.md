## 安装运行
安装
pip install jupyter

运行远程密码访问

jupyter notebook --generate-config 生成配置

jupyter notebook password 设置密码

修改config中的 

- server.ip 为 0.0.0.0 
- server.port 
- server.allow_remote..
- server.allow_root..


运行
jupyer notebook