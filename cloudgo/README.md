# go 服务端程序 cloudgo
## 程序功能
用户使用GET方法向服务器端发送数据：
```
http://localhost:8080/?xiaoming=?
```
服务器端根据用户发送的姓名数据，查询对应的学号，并返回给用户。

---
## 实现思路
### 获取用户数据
获取用户使用 ```GET``` 方法提交的数据：
```go
mx.HandleFunc("/", testHandler(formatter)).Methods("GET")
```
其中，通过 ```testHandler``` 函数处理本次连接：
```go
func testHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        // TODO
    }
}
```
### 读取数据
程序读取 json 文件，使用 ```encoding/json``` 包解析成```map```型的键值对：
```go
func readFile(filename string) (map[string]string, error) {
    var data map[string]string
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Println("ReadFile: ", err.Error())
        return nil, err
    }
    err = json.Unmarshal(bytes, &data)
    if err != nil {
        fmt.Println("Unmarshal: ", err.Error())
        return nil, err
    }

    return data, nil
}
```
---
## curl 测试
```datadata.json``` 中的数据如下：
```json
{
  "wxl": "16340241",
  "xiaoming": "16340001",
  "trump": "16340002"
}
```
**测试结果**
```
$ curl http://localhost:8080/?wxl=?
16340241

$ curl  http://localhost:8080/?xiaoming=?
16340001

$ curl  http://localhost:8080/?trump=?
16340002

$ curl  http://localhost:8080/?nobody=?  ### 姓名不存在
This name does not exist!
```
---
## ab 测试
执行 ab 测试：
```
$ ab -n 1000 -c 100 http://localhost:8080/?wxl=?
This is ApacheBench, Version 2.3 <$Revision: 1430300 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /?wxl=?                       ###请求的资源
Document Length:        9 bytes                       ###文档返回的长度，不包括相应头

Concurrency Level:      100                           ###并发个数
Time taken for tests:   0.172 seconds                 ###总请求时间
Complete requests:      1000                          ###总请求数  
Failed requests:        0                             ###失败的请求数
Write errors:           0
Total transferred:      125000 bytes
HTML transferred:       9000 bytes
Requests per second:    5799.45 [#/sec] (mean)        ###平均每秒的请求数
Time per request:       17.243 [ms] (mean)            ###平均每个请求消耗的时间
Time per request:       0.172 [ms] (mean, across all concurrent requests)
Transfer rate:          707.94 [Kbytes/sec] received  ###传输速率

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    4   2.2      4       9
Processing:     1   12   5.9     11      31
Waiting:        1   10   5.5      9      28
Total:          4   16   5.3     15      32

Percentage of the requests served within a certain time (ms)
  50%     15                                          ###50%的请求都在15ms内完成
  66%     18
  75%     20
  80%     21
  90%     23
  95%     27
  98%     30
  99%     31
 100%     32 (longest request)
```
