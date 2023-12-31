# instant-messaging-system-based-on-Golang
#### 说明
* 刚进入大二的暑假，在同学见一下报名参加字节跳动的后端青训营，在学了一些知识后按照教程用Golang做了一个即时通讯系统，从v1到v8一步步改进

#### 版本说明
* v1版本：实现了tcp链接的建立
* v2版本：在服务器端加入了用户类型
* v3版本：增加了用户消息的全体广播功能
* v4版本：增加了在线用户查询的功能
* v5版本：增加了修改用户名功能
* v6版本：增加了对用户的超时强踢功能
* v7版本：增加了私聊功能
* v8版本：创建并完善了客户端内容

#### 运行方法
* 以v8以前的版本
  * 进入文件夹后运行`go build`命令，然后运行生成的文件即可启动服务器
  * 对于用户，在新开的终端中输入`nc 127.0.0.1 8888`
* 对v8版本
  * 服务器端端启动和之前一样
  * 客户端的启动方式是，进入`client`文件夹后在终端中输入`go run client.go`
