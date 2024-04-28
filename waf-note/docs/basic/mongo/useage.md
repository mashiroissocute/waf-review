## 连接
./mongo ip:port/DB  -u $user -p 'xxxxxxx' 

## 使用
展示db ： show dbs

使用db： use xxx

展示集合： show collections

展示文档：db.getCollection("same-domain-list").find()

插入文档 ： db.getCollection("same-domain-list").insert({"appid" : 251011484, "status": 1})

https://www.runoob.com/mongodb/mongodb-tutorial.html


## 不能使用GORM操作MONGO的原因
In short: you can't. [GORM](http://gorm.io/) is created for relational databases, and MongoDB is not a relational but a NoSQL database.
And you can't even use GORM with all SQL databases, the [officially supported list](http://gorm.io/docs/connecting_to_the_database.html#Supported-Databases) at the moment is: MySQL, PostgreSQL, SQLite3 and SQL Server, although you can "easily" add support for other SQL servers by [writing GORM dialects](http://gorm.io/docs/dialects.html) for them. But that's the end of it. Adding support for MongoDB would require more work than what your gain would be.
Consider using the [official MongoDB driver](https://github.com/mongodb/mongo-go-driver) which is quite mature now. Or if using GORM is a must for you, you must choose another (not MongoDB, preferably one of the above listed supported) database.