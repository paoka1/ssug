## SSUG

Shauio's short URL generator，短链接服务API，使用Golang实现，使用SQLite持久化数据

### 如何使用

在Releases中下载对应架构的最新版本或自己构建，命令行参数：

1. -host：监听的地址，默认127.0.0.1
2. -key：访问管理接口的密钥，默认key123456
3. -len：生成的短链最小长度，默认3
4. -port：监听的端口，默认8000
5. -ttl：短链存活的时长（秒），默认24小时
6. -debug：开启debug模式

### 使用API

API 列表如下：

1. 添加短链：`host:port`+`/add`，POST参数，`key`为访问密钥、`url`为要添加的短链的原始链接
2. 访问短链：`host:port`+`/短链`

使用Python编写的一个[API实例](example/api.py)
