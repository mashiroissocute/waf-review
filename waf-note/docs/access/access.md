## 0.5 saaswaf

### cvm时代
-	集群包含预先分配的公共集群池 + 申请的独立集群
-	所有的独享集群都是云梯集群
-	云梯支持gz sh cd hk bj sing接入， 目前购买saas实例的地域支持gz sh cd hk bj。意味着新购的独享版实例都将申请云梯集群
-	集群级别（level） 0 中小客户共享集群 1大客户共享集群 2独立集群
-	集群类别（type） 0 ziyan 1 yunti
-	集群版本 （package id） 4 高级版 5 企业版 6 旗舰版 7 标准版 8 独享版 9 小微企业版
-	云梯集群的版本只有小微/独享/其他
-	新接入实例（不存在同地已接入实例）都是分配的云梯（独享或者公共）


- 集群选择：
	- instance已接入集群，（instance下已经存在接入的域名）在t_spart_instance中根据instance_id查找group_id即可（可能是ziyan也可能是yunti）。
	- instance未接入集群，（instance下的第一条域名接入）：
		-  存在同地已接入的实例，instance_level != 6 非独享版实例，分配用户名下相同地域其他instance所在的集群（ziyan、yunti）（用户有多实例，且其他同地域实例中已有域名接入）
		-  不存在同地已接入的实例，instance_level = 6 独享版本实例，申请新云梯集群（yunti）
		-  不存在同地已接入的实例，instance_level !=6 非独享版实例，根据实例信息查找yunti公共集群（yunti）  
		- 找到groupid后，向t_sparta_instance中插入一条该instance的数据，以后该instance变为已接入集群状态


- 新增集群（独享版 yunti）：
	- 可用区（Zone）是指腾讯云在同一地域内电力和网络互相独立的物理数据中心。其目标是能够保证可用区间故障相互隔离（大型灾害或者大型电力故障除外），不出现故障扩散，使得用户的业务持续在线服务。通过启动独立可用区内的实例，用户可以保护应用程序不受单一位置故障的影响。

	- t_sparta_group添加一条数据得到groupid 。
	- 分配cvm机器: 在t_sparta_nginx_ip 中选择groupid为0且地域满足条件的机器，根据策略尽量均匀的分配在不同的可用区。
	- 修改cvm机器的集群id，将其挂载到groupid下。 