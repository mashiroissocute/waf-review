## kafka生产者 : log-udp-server
https://www.cnblogs.com/lys_013/p/11957479.html
设置KAFKA生产者属性，并返回一个KafkaProducer
``` 
import java.util.Properties

import org.apache.kafka.clients.producer.KafkaProducer

object kafkaProducer {
# 一级标题
  def genKafkaProducer(): KafkaProducer[String, String] = {
      val props = new Properties()
      props.put("bootstrap.servers", brokers) //"100.119.167.50:6328"
      props.put("acks","all") //acks：指定必须要有多少个分区的副本接收到该消息，服务端才会向生产者发送响应，可选值为：0,1,2，…，all（-1），如果设置为0，producter就只管发出不管kafka server有没有确认收到。设置all则表示kafka所有的分区副本全部确认接收到才返回。
	  
      val a :java.lang.Integer = 0
      var b :java.lang.Integer = 16*1024
      val c :java.lang.Integer = 1
      val d :java.lang.Integer = 33554432
	  
      props.put("retries", a) //生产者向kafka发送消息可能会发生错误，有的是临时性的错误，比如网络突然阻塞了一会儿，有的不是临时的错误，比如“消息太大了”，对于出现的临时错误，可以通过重试机制来重新发送。这个参数还是非常重要的，在生产环境中是必须设置的参数，为设置消息重发的次数。在Kafka中可能会遇到各种各样的异常（可以直接跳到下方的补充异常类型），但是无论是遇到哪种异常，消息发送此时都出现了问题，特别是网络突然出现问题，但是集群不可能每次出现异常都抛出，可能在下一秒网络就恢复了呢，所以我们要设置重试机制。
	  
	  props.put("retry.backoff.ms", 100)//设置隔多久重试一次
	  
	props.put("buffer.memory", d) //生产者的数据是先放在缓冲区的，同时会有一个独立线程Sender去把消息分批次包装成一个个Batch。生产者的内存缓冲区大小。如果生产者发送消息的速度 > 消息发送到kafka的速度，那么消息就会在缓冲区堆积，导致缓冲区不足。这个时候，send()方法要么阻塞，要么抛出异常。所以，如果buffer.memory设置的太小，可能导致的问题是：消息快速的写入内存缓冲里，但Sender线程来不及把Request发送到Kafka服务器，会造成内存缓冲很快就被写满。而一旦被写满，就会阻塞用户线程，不让继续往Kafka写消息了。
	
      props.put("batch.size", b) //生产者在发送消息时，可以将即将发往同一个分区的消息放在一个批次里，然后将这个批次整体进行发送，这样可以节约网络带宽，提升性能。该参数就是用来规约一个批次的大小的。如果设置的太小，会造成消息发送次数增加，会有额外的IO开销。如果把该参数调的比较大的话，也是不会造成消息发送延迟的，只会占用比较大的内存，因为还有参数linger设置固定多长时间，就算没塞满Batch，也会发送。批次的大小默认是16K，这里设置了32K，设置大一点可以稍微提高一下吞吐量，设置这个批次的大小还和消息的大小有关，假设一条消息的大小为16K，一个批次也是16K，这样的话批次就失去意义了。所以我们要事先估算一下集群中消息的大小，正常来说都会设置几倍的大小。
	  
      props.put("linger.ms", c) //生产者在发送一个批次之前，可以适当的等一小会，这样可以让更多的消息加入到该批次。这样会造成延时增加，但是降低了IO开销，增加了吞吐量。
	  

	  
      props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer")//key.serializer：关键字的序列化方式
      props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer")//value.serializer：消息值的序列化方式
      
      val producer = new KafkaProducer[String,String](props)
      producer

  }

}
```
### 属性设置
#### acks
0	 Producer 往集群发送数据不需要等到集群的返回，不确保消息发送成功。安全性最低但是效率最高。
1	 Producer 往集群发送数据只要 Leader 应答就可以发送下一条，只确保 Leader 接收成功。Leader接受后还未同步给replica的时候挂掉，会导致消息丢失。
-1 或 all	 Producer 往集群发送数据需要所有的ISR Follower 都完成从 Leader 的同步才会发送下一条，确保 Leader 发送成功和所有的副本都成功接收。安全性最高，但是效率最低。



## 生产数据
``` 
访问日志topic : filtered_access_log
攻击日志topic : con_attack_log
错误日志topic : report
心跳日志topic : heart
请求经历模块日志 : strace

    val address = InetAddress.getByName("127.0.0.1")
    val socket = new DatagramSocket(port, address) //端口
    val mainProducer = kafka.kafkaProducer.genKafkaProducer() //返回KafkaProducer

    maxMessageLength=16*1024
    
    val bytes = new Array[Byte](maxMessageLength)
    val packet = new DatagramPacket(bytes, bytes.length)
    while (true) {
      socket.receive(packet) //读取udp报文（内容是json）
      val sd = new String(packet.getData, 0, packet.getLength, "UTF-8")
      val rec1 = new ProducerRecord[String, String](topic, sd) //topic,msg（json）
      mainProducer.send(rec1)
    }
```