# Charging Station

本项目基于 [kratos框架](https://github.com/go-kratos/kratos) 搭建基础微服务项目，开发各种独立服务,bff层作为入口，整合各个服务。


#### 管理后台服务
- [x]  管理后台

#### 管理员服务
- [x]  登录退出
- [x]  管理员管理

#### 权限服务
- [x]  角色管理
- [x]  菜单管理
- [x]  权限管理
- [x]  api管理

#### 文件服务
- [x]  前端直传(获取token)

#### 运行方式

##### 组件
管理后台web框架: vue-element-admin

服务注册与发现： ETCD

链路追踪：jaeger

数据库：mysql

缓存：redis

orm: GORM



### 运行(以权限服务为例)
##### 数据库
1.导入sql
```
文件位于initSql/authorization.sql
```

##### 后端
1.下载
```
$ go clone git@github.com:ZQCard/authorization.git
```

2.安装依赖
```
$ cd authorization && go mod tidy
```

3.设置配置 以管理员服务为例, 配置文件位于configs/
```
$ vim ./app/administrator/configs/config.yaml
```

4.运行项目
```
$ kratos run
```

5.如果使用docker部署,请更改Dockerfile的启动配置文件 config-dev.yaml 并更改Makefile中的docker命令为自己的配置
```
$ make docker
```

##### [管理后台前端](https://repo.example.com/frontend)
1.安装依赖
```
$ cd web && npm install
```

2.启动项目
```
$ npm run dev
```

### 部署(docker)
##### 后端
可以参考kratos部署 (https://go-kratos.dev/docs/devops/docker)

1.服务部署 以管理员服务为例 app/
```
$ cd app/administrator
```
2.make打包docker镜像
```
# PS:如果是打包admin镜像 app/project/admin 请执行 make dockerAdmin
$ make docker
```
3.运行容器 
```
# 注意端口映射设置， docker部署容器9000端口， 本地开发端口不能全是9000
docker run -p 9000:9000 --name kratos-base-project-administrator --restart=always -v /data/project/kratos-base-project/app/administrator/configs:/data/conf -d kratos-base-project/administrator:0.1.0
```
##### 前端
1.进入前端目录
```
$ cd web
```
2.编译
```
$ npm run build:prod
```
3.将dist文件夹上传至服务器

##### nginx示例
```
server {
  listen 80;
  listen [::]:80;
  server_name kratos.niu12.com;
  index index.html;
  root /data/project/kratos-base-project/web/dist;

  # 管理后台接口转发代理
  location  /api/ {
  # nginx代理设置Header
      proxy_set_header            X-real-ip $remote_addr;

      proxy_pass                  http://127.0.0.1:8000/;
  }

}
```

* 有任何建议，请扫码添加我微信进行交流。

![扫码提建议](https://kratos-base-project.oss-cn-hangzhou.aliyuncs.com/f8f5dacdf87cf358c98c9eb60ce2a13.jpg)
