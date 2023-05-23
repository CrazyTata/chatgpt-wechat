[TOC]
## 部署相关
### 部署chat
```shell
cd /www/wwwroot/chatgpt-wechat/chat/  #进入项目目录
git pull #拉取最新代码
vi ./service/chat/api/etc/chat-api.yaml #修改配置文件，如果配置无变动，跳过这一步
sudo docker-compose build #打包
sudo docker-compose down #关闭服务
sudo docker-compose up -d #启用服务
docker logs tail=100 chat-web-1 #查看接口日志
```

### 部署script
```shell
cd /www/wwwroot/chatgpt-wechat/script/  #进入项目目录
git pull #拉取最新代码
vi ./script/script-api.yaml #修改配置文件，如果配置无变动，跳过这一步
go build #打包
./script & #启用服务
```

### 部署cron
```shell
cd /www/wwwroot/chatgpt-wechat/cron/  #进入项目目录
git pull #拉取最新代码
vi ./cron/cron-api.yaml #修改配置文件，如果配置无变动，跳过这一步
go build #打包
./cron & #启用服务
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
### 将 fork 的仓库合并到自己的仓库
#### 在本地电脑上，将自己仓库的代码克隆到本地：
```
git clone git@github.com:<your-username>/<a>.git
git@github.com:chy4pro/chatgpt-wechat.git

cd <a>
git remote add upstream git://github.com/<original-author>/<b>.git
git remote add upstream git@github.com:whyiyhw/chatgpt-wechat.git


git fetch upstream
git merge upstream/main


git push origin main

```