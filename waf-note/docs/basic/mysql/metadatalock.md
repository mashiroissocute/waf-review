## 一、MDL背景
在MySQL中，元数据锁（Metadata Lock）是一种用于控制对数据库对象（如表、存储过程和触发器）访问的锁。当一个对象被获取了元数据锁，它会阻止其他事务修改或删除该对象，直到锁被释放。
当事务执行某些操作时，如创建或删除表，或更改表的结构，MySQL会要求获取元数据锁。从而确保数据一致性，防止事务间冲突。
申请MDL锁的事务由队列维护。因此，如果一个事务长时间持有元数据锁，它会导致后续事务等待。在业务高峰期，大量事务提交的情况下，容易导致数据库连接堆积，进而导致cpu飙升，服务器宕机。

## 二、MDL处理
当线上MySQL出现大量MDL问题时，应该及时找到导致阻塞的事务ID，执行KILL，防止服务器宕机。示例如下：

- 1执行DDL:
```sql
ALTER TABLE  
    tb_instance_conf  
ADD COLUMN 
    api_security  int(2)  DEFAULT '0' COMMENT 'api安全开关';
```

- 2导致MDL堆积：

此时MDL锁队列已经有大量的事务堆积，即使是查询语句也被阻塞。
线上业务查询不到数据，影响用户控制台使用。

- 3Kill：

找到引起MDL堆积的DDL语句，执行Kill，放弃执行。
Kill 452083692

MDL堆积消失，线上业务恢复。

## 三、MDL预防
为了预发DDL操作带来的MDL堆积问题，以下有一些提前检测的方法：

1）假如表上存在正在运行的长事物，阻塞DDL，继而阻塞所有同表的后续操作。
通过以下语句查询，表上正在运行的事务。等到事务执行完毕后，再做DDL操作。

```sql
select * from information_schema.processlist where COMMAND != "Sleep"
```

2）加入表上存在未提交事物，阻塞DDL，继而阻塞所有同表的后续操作。
通过场景一中的sql，看不到表上有任何操作。
我们通过以下语句，再查看表上是否有未提交的事物。
```sql
select * from information_schema.innodb_trx
```

事务没有提交，表上的锁不会释放。因此，我们等事务提交后，再做DDL操作。

## 四、DDL变更
通过MDL预防检测之后，再进行DDL变更，通常可以顺利的执行。
但是，并非绝对。
因为，我们只能通过上述方式减少发生MDL的概率，却无法完全避免问题。
那么，如何安全的进行DDL变更呢？
一般来说，热点表的DDL变更，有较大概率导致MDL带来的服务器高负载问题。
安全的进行热点表DDL有以下三个方法：

- 分表
将大表，划分为多个小表。分表又分为横向切分和纵向切分。
横向切分：将总表数据分散到多个子表。根据数据量的规模来划分，保证单表的容量不会太大，从而减少单表MDL概率。
纵向切分：将总表的多个列，根据活跃度进行分离。把DDL放到低活跃度的子表上。
分表的缺点是，需要进行大量的业务改造。

- Pt工具
pt-online-schema-change 可以在线修改表结构。主要原理为：
1.创建一个和要执行 alter 操作的表一样的新的空表结构(是alter之前的结构) 
2.在新表执行alter table 语句（速度应该很快） 
3.在原表中创建触发器3个触发器分别对应insert,update,delete操作 
4.以一定块大小从原表拷贝数据到临时表，拷贝过程中通过原表上的触发器在原表进行的写操作都会更新到新建的临时表 
5.Rename 原表到old表中，在把临时表Rename为原表 
Pt工具的缺点是，Rename需要获取MDL，所以Pt工具仍然可能遭遇MDL堆积问题。并且Pt工具不允许表上已经存在触发器。


- 主从切换（采纳）
在MySQL从节点上进行DDL，后进行主从切换。
相比主节点，从节点仅提供读的能力，在进行DDL的时候。发生MDL堆积的概率比主节点更小。因此，我们可以在从节点上进行变更后，进行主从切换。


参考：

[MySQL出现Waiting for table metadata lock的原因以及解决方法](https://www.cnblogs.com/digdeep/p/4892953.html)

[MySQL之使用pt-online-schema-change在线修改大表结构](https://blog.51cto.com/u_11045899/5968036)



## 五、补充
关于DDL DML DCL

- DDL（Data Definition Languages）语句：数据定义语言，这些语句定义了不同的数据段、数据库、表、列、索引等数据库对象的定义。常用的语句关键字主要包括 create、drop、alter等。
- DML（Data Manipulation Language）语句：数据操纵语句，用于添加、删除、更新和查询数据库记录，并检查数据完整性，常用的语句关键字主要包括 insert、delete、udpate 和select 等。(增添改查）
- DCL（Data Control Language）语句：数据控制语句，用于控制不同数据段直接的许可和访问级别的语句。这些语句定义了数据库、表、字段、用户的访问权限和安全级别。主要的语句关键字包括 grant、revoke 等。

