# 极速百家乐游戏服go实现

极速百家乐游戏服go实现

### 目录结构
```go
├─classic
│  ├─bet_handler
│  ├─draw_hander
│  ├─event
│  ├─game_type
│  ├─order
│  │  └─cache
│  ├─sign_handler
│  └─ws_message_handler
├─dragon_tiger
│  ├─bettor
│  ├─db
│  ├─drawer
│  ├─event
│  └─sign
├─framework
│  ├─async
│  ├─base
│  ├─env_config
│  │  ├─dev
│  │  └─uat
│  ├─game_server
│  │  ├─conf
│  │  │  └─snow_flake_id
│  │  ├─error_code
│  │  ├─game
│  │  │  ├─controller
│  │  │  ├─dao
│  │  │  │  ├─gamedb
│  │  │  │  ├─redisdb
│  │  │  │  └─uiddb
│  │  │  └─service
│  │  │      ├─base
│  │  │      ├─bet
│  │  │      │  └─bet_dto
│  │  │      ├─draw
│  │  │      ├─game
│  │  │      ├─game_event
│  │  │      ├─interface
│  │  │      │  ├─bet
│  │  │      │  ├─dao
│  │  │      │  ├─draw
│  │  │      │  ├─events
│  │  │      │  └─ws_message
│  │  │      ├─listenner
│  │  │      ├─sign
│  │  │      │  ├─impl
│  │  │      │  ├─interfaces
│  │  │      │  └─sign_dto
│  │  │      ├─type
│  │  │      │  ├─dto
│  │  │      │  └─VO
│  │  │      └─websocket_message
│  │  ├─mq
│  │  ├─redis
│  │  │  ├─cache
│  │  │  ├─rediskey
│  │  │  ├─redis_tool
│  │  │  └─types
│  │  └─rpc_client
│  │      └─dto
│  ├─network
│  ├─resource
│  │  ├─base_socket
│  │  ├─cache
│  │  ├─conf
│  │  │  └─observer
│  │  ├─console
│  │  ├─db
│  │  │  └─manager
│  │  ├─error
│  │  ├─httpserver
│  │  │  └─beego
│  │  │      ├─controllers
│  │  │      ├─log
│  │  │      └─middlewares
│  │  ├─nacos
│  │  ├─protocol
│  │  ├─snow_flake_id
│  │  └─uiddb
│  ├─tool
│  └─trace
│      ├─alils
│      └─es
└─speed
    ├─bet_handler
    ├─draw_handler
    ├─event
    ├─game_type
    ├─order
    │  └─cache
    ├─sign_handler
    └─ws_message_handler
```