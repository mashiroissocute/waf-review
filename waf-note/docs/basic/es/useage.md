## 安装

- ES
非root启动
https://blog.csdn.net/qq_32690999/article/details/78505533
- Kibana
--alow-root

## RestfulAPI使用
https://endymecy.gitbooks.io/elasticsearch-guide-chinese/content/getting-started/modifying-data.html

## url
`url -X<REST Verb> <Node>:<Port>/<Index>/<Type>/<ID>`
``` 
curl -XPOST 'localhost:9200/kibana_sample_data_flights/_search?pretty' -d '{"query": { "match_all": {} }}' -H 'Content-Type: application/json' -iv

{
    "_index" : "kibana_sample_data_flights",
    "_type" : "_doc",
    "_id" : "hpdeRYMBHHvZZwP9GN2K",
    "_score" : 1.0,
    "_source" : {
        "FlightNum" : "SNI3M1Z",
        "DestCountry" : "IT",
        ...
    }
}
```
## 查询节点信息
## 查询/设置/删除索引index信息
## 查询记录信息

## java client
https://endymecy.gitbooks.io/elasticsearch-guide-chinese/content/java-api/client.html

``` 
public static void main( String[] args ) throws IOException
    {
        //client
        RestHighLevelClient client = new RestHighLevelClient(RestClient.builder(new HttpHost("localhost", 9200)));

        //add source
        String index = "test1";
        String type = "_doc";
        // 唯一编号
        String id = "1";
        IndexRequest request = new IndexRequest(index, type, id);
        Map<String, Object> jsonMap = new HashMap<>();
        jsonMap.put("uid", 1234);
        jsonMap.put("phone", 12345678909L);
        jsonMap.put("msgcode", 1);
        jsonMap.put("sendtime", "2019-03-14 01:57:04");
        jsonMap.put("message", "xuwujing study Elasticsearch");
        request.source(jsonMap);
        IndexResponse indexResponse = client.index(request);
    }
```


``` 
PUT /test1/_doc/1?timeout=1m HTTP/1.1
Content-Length: 118
Content-Type: application/json
Host: localhost:9200
Connection: Keep-Alive
User-Agent: Apache-HttpAsyncClient/4.1.2 (Java/1.8.0_131)

{"msgcode":1,"uid":1234,"phone":12345678909,"sendtime":"2019-03-14 01:57:04","message":"xuwujing study Elasticsearch"}

HTTP/1.1 200 OK
content-type: application/json; charset=UTF-8
content-length: 153

{"_index":"test1","_type":"_doc","_id":"1","_version":3,"result":"updated","_shards":{"total":2,"successful":1,"failed":0},"_seq_no":2,"_primary_term":1}
```


