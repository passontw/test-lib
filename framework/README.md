# 游戏服通用模块

## 使用方式

公共模块module名称为sl.framework.com

使用git submodule add 命令放入本地，统一使用文件名称 `framework`

```bash
git submodule add http://gitlabdev.solidleisure.com:5999/g32-group/g32-common-game-server.git framework
```

```go
//修改go.mod 增加如下一行 将module指向本地模块
replace sl.framework.com => ./framework
```

## 游戏服部署以及配置说明

- 新服务部署之前务必通知到中台 前端 产品 运营

- RocketMQ
    - 配置服务所需要的rocketmq topic(所有游戏服上线，创建一次即可，需要检查服务对应的topic是否存在)

- SQL
    - 通知运维创建所需要的数据库
    - 将创建表的SQL上传到以下confluence中 https://wiki.slleisure.com/pages/viewpage.action?pageId=853653
    - 通知运维执行SQL

- Nacos配置
    - 在以下网址 http://10.146.40.240:30011 创建需要服务的Nacos配置
    - Nacos配置需要在所有环境相关命名空间下都创建，G32环境对应的开发、测试、UAT、Prod的命名空间为：g32-dev、g32-qa、g32-uat、g32-prod
    - 配置要能明显区分是本服务的配置，Nacos配置以yml形式配置，配置文件已yml为后缀名，如roulette-game-server.yml
    - Nacos创建完毕后通知运维，由运维将配置转移到prod环境地址 http://nacos.g32-prod.com，线上环境，游戏服配置命名空间为：G32-Game

- 流水线部署
    - 在KubeSphere部署流水线，并运行流水线成功
    - 通知运维在如下地址中创建Jenkins job https://uat-jenkins.slleisure.com/
    - 通知运维在UAT环境创建对应游戏服服务，如有变量需要注入，则与运维沟通，注入相应变量
    - 在以下网址观察服务是否正常启动以及日志是否正常 https://ppu-kubesphere.slleisure.com
    - 在 https://ppu-kubesphere.slleisure.com 中拷贝服务dns，发给中台，由中台做相应配置，或者在 总控-游戏管理-游戏配置 中修改

- 配置备份
    - 将各个环境Nacos配置上传到本服务代码仓库
    - 将各个环境Jenkinsfile配置上传到本服务代码仓库