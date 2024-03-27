## Yarn
在MapReduce中，可以看出MapReduce既作为了计算框架，有作为了资源调度框架。耦合之后，不利于扩展。
使用Yarn作为资源调度框架，MapReduce作为计算框架的结构更为清晰。

## 架构
![alt text](image-7.png)
![alt text](image-8.png)

MapReduce如果想在Yarn上运行，需要开发遵循Yarn规范的MapReduce ApplicationMaster。