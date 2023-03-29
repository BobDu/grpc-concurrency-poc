# grpc-concurrency-poc
验证 grpc 的并发使用方式


server

```shell
cd grpc-server
go build server.go && ./server
```

web-server (also grpc-client)
```shell
cd go-web
go build -o web && ./web
```


测试流程
===

1. 不启动 grpc-server 只启动 web-server
```
此时 web-server 可以成功启动 并没有因为 grpc无法访问而阻塞影响启动
curl http://127.0.0.1:50052/a
{"resp":""}
日志
[GIN-debug] GET    /:name                    --> main.main.func1 (3 handlers)
   rpc err rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing: dial tcp 127.0.0.1:50051: connect: connection refused"

[GIN] 2023/03/29 - 23:12:49 | 200 |     328.641µs |       127.0.0.1 | GET      "/a"
```

2. 此时再启动 grpc-server
```
curl http://127.0.0.1:50052/a
{"resp":"Hey, a!"}
日志
Hey, a!
[GIN] 2023/03/29 - 23:15:05 | 200 |    7.058257ms |       127.0.0.1 | GET      "/a"
```

3. 再多次启动& 关闭 grpc-server 并多次 curl请求 web-server
发现只要 grpc-server处于正常状态就可以正常响应 否则会快速失败

该POC验证了
===

1. grpc客户端Dial()在服务器未正常工作情况下创建的client 再后续服务器恢复后可以仍可自动重新建立链接正常响应
2. grpc client 的rpc方法在服务器未正常工作时会快速失败 返回错误
rpc err rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing: dial tcp 127.0.0.1:50051: connect: connection refused"
