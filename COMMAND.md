## 部署相关
### 部署代码步骤
```shell
cd /www/wwwroot/chatgpt-wechat/chat/  #进入项目目录
git pull #拉取最新代码
vi ./service/chat/api/etc/chat-api.yaml #修改配置文件，如果配置无变动，跳过这一步
sudo docker-compose build #打包
sudo docker-compose down #关闭服务
sudo docker-compose up -d #启用服务
docker logs tail=100 chat-web-1 #查看接口日志
```

### docker
清理磁盘 `docker system prune`

## 框架相关
### 生成model
```
cd service/user/model
goctl model mysql ddl -src user.sql -dir . -c
```

### 生成api
```
cd book/service/user/api
goctl api go -api user.api -dir .
```

## git相关