package types

import (
	"sync"
	"time"
)

type (
	// HashRedisInfo hash的redis信息
	HashRedisInfo struct {
		HTable string        //Hash表名
		Filed  string        //Hash表中的Filed字段
		Expire time.Duration //Hash表的过期时间
	}

	RedisInfo struct {
		Key    string
		Expire time.Duration
	}

	RedisLockInfo struct {
		*RedisInfo
		Owner int64 //key持有者
	}
)

const (
	BacDuration              = time.Duration(15) * time.Second //电子百家乐为最短时长游戏
	GameRoundDurationDefault = time.Duration(30) * time.Second //默认局时长
	RedisLockExpireDuration  = time.Duration(2) * time.Second  //redis lock过期时间

)

var (
	RedisExpireDurationRWMutex sync.RWMutex
	RedisExpireDuration        = GameRoundDurationDefault * 2

	ServerRedisKeyPrefix = "" //服务公共前缀 初始化为空
)
