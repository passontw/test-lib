package dao

import (
	"sl.framework.com/game_server/conf"
	"sl.framework.com/trace"
	"testing"
)

func TestOrmInit(t *testing.T) {
	_ = trace.LoggerInit()
	conf.NacosClientInitOnce(conf.OnLocalConf)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "name1", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := OrmInit(); (err != nil) != tt.wantErr {
				t.Errorf("OrmInit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
