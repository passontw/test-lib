package conf

import "testing"

func Test_getDetailInHost(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 uint64
	}{
		{name: "test1", args: args{host: "http://nacos.nacos:8848"}, want: "nacos.nacos", want1: 8848},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getDetailInHost(tt.args.host)
			if got != tt.want {
				t.Errorf("getDetailInHost() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getDetailInHost() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
