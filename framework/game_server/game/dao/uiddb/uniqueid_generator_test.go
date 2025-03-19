package uiddb

import (
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/game/dao/gamedb"
	"sl.framework.com/tool"
	"sl.framework.com/trace"
	"testing"
	"time"
)

func TestOrmUidDbInit(t *testing.T) {
	_ = trace.LoggerInit()
	conf.NacosClientInitOnce(conf.OnServerConf)
	_ = gamedb.OrmGameDbInit()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "name1", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := OrmUidDbInit(); (err != nil) != tt.wantErr {
				t.Errorf("OrmUidDbInit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	g := &GoWorkerNode{}
	g.QueryMaxId()
	g.Id = g.Id + 1
	g.CreateTime = time.Now()
	g.HostName = tool.GetLocalIp()
	g.Insert()
}
