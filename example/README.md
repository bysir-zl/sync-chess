# 四川麻将

> 源代码量900行左右. 主要代码量集中在玩家出牌逻辑处理, 其实这都是些低逻辑代码, 实际你要写的代码更少.


## 依赖

### 网络库
使用ws协议通讯, 一个请求一个协程对应一个handle, 清晰明了
github.com/bysir-zl/hubs

## 运行
启动服务器
```
go run /server/main.go
```
启动客户端

随便找个浏览器打开client.html即可, 然后在地址后面输入 xxxx/client.html **#1** 就代表玩家id1

依次打开 #2 #3, 这时候系统将会发牌, 开始游戏 .


## 期望架构
推荐一个简单的集群架构
```
// ->>> 代表可有多个服务
client -> agent -> db -(find chess server and conn)>>> chess server ->>> log server
```
