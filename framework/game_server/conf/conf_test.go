package conf

import (
	"log"
	"testing"
)

func Test_parseConfLine(t *testing.T) {
	type args struct {
		line string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test1", args: args{"${MQ_NAMESERVER:http://rocketmq-nameserver-dev.mq:9876}"}, want: "test1"},
		{name: "test2", args: args{"${MQ_NAMESERVER}"}, want: "test2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Println(tt.name, ":", parseConfLine(tt.args.line))
		})
	}
}
