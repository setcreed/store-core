# store-core
使用grpc的方式提供仓储层crud的服务

# 使用理念
作为一个提供仓储层crud的服务，外部服务只需要调用该服务接口即可，通过修改配置文件中的sql语句做增删改查功能。

# 使用方式
准备一个mysql

修改`config.yaml`
```yaml
default:
  mode: debug
  app:
    rpcPort: 8080
    httpPort: 8090

dbConfig:
  dsn: "root:123456@tcp(127.0.0.1:33060)/test?charset=utf8mb4&parseTime=True&loc=Local"
  maxOpenConn: 20
  maxLifeTime: 1800
  maxIdleConn: 5
```

## 启动
```bash
make run
```
![](docs/img/run.png)


## 查询sql

在`config.yaml`中配置 sql语句
```yaml
sqlConfig:
  - name: userList
    sql: "select * from users where user_id>@id"
```
可以自己构造查询方法，可以仿照: [examples/client/client.go](examples/client/client.go)

也可以使用[store-core-sdk](https://github.com/setcreed/store-core-sdk)项目,  例子：https://github.com/setcreed/store-core-sdk/blob/master/examples/query.go


只要在`config.yaml`中配置sql语句，就可以使用

支持配置重载：
```bash

./store-core reload --configfile=config.yaml

{"function_name":"github.com/setcreed/store-core/cmd/app.NewServerCommand.func4","level":"info","line_num":81,"module_name":"/Users/sss/workspace/WickCloud/store-core/cmd/app/server.go","msg":"[配置文件重载成功]","time":"20T19:31:42+08:00"}

```
