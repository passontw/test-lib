package gamedb

import (
	"sl.framework.com/game_server/conf"
	"sl.framework.com/game_server/game/dao/uiddb"
	"sl.framework.com/trace"
	"testing"
)

func TestOrmGameDbInit(t *testing.T) {
	_ = trace.LoggerInit()
	conf.NacosClientInitOnce(conf.OnLocalConf)

	if err := uiddb.OrmUidDbInit(); err != nil {
		trace.Error("TestOrmGameDbInit uiddb OrmUidDbInit error=%v", err.Error())
		return
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "name1", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := OrmGameDbInit(); (err != nil) != tt.wantErr {
				t.Errorf("OrmGameDbInit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
