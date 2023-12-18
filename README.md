## SSUG

Shauio's short URL generator，短链接服务API，使用Golang实现，使用SQLite持久化数据

### 使用SSUG

在[Releases](https://github.com/paoka1/ssug/releases)中下载对应架构的最新版本或自己构建，命令行参数：

1. -host：监听的地址，默认127.0.0.1
2. -key：访问管理接口的密钥，默认key123456
3. -len：生成的短链最小长度，默认3
4. -port：监听的端口，默认8000
5. -ttl：短链存活的时长（秒），默认24小时
6. -debug：开启debug模式

### 使用API

API 列表如下：

1. 为原始链接添加短链：

   ```api
   路径：host:port/add
   方法：POST
   参数：
   	key：访问密钥
   	url：原始链接
   返回值：
   	json格式：
   		code：状态码
   		msg：错误信息
   		data：返回的短链
   ```

   注意：添加已存在的链接时，返回code不为0，data为存在链接的短链

2. 使用短链访问原始链接：

   ```api
   路径：host:port/{短链}
   方法：GET
   参数：
   	URL路径参数：{短链}
   返回值：
   	成功返回javascript，进行跳转
   	失败返回json格式：
   		code：状态码
   		msg：错误信息
   		data：空
   ```

使用Python调用API的一个例子：[api.py](example/api.py)
