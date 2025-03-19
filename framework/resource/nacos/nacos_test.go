package nacos

import (
	"path/filepath"
	"testing"
)

func Test_extractHostAndPort(t *testing.T) {
	type args struct {
		urlStr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"x", args{urlStr: "10.146.40.110:8848"}, "10.146.40.110:8848", false},
		{"y", args{urlStr: "http://nacos.nacos:8848"}, "nacos.nacos:8848", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractHostAndPort(tt.args.urlStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractHostAndPort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractHostAndPort() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filePath(t *testing.T) {
	type args struct {
		path        string
		defaultPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"xxx",
			args{
				path:        localFilePath("config.yaml"),
				defaultPath: filepath.Join(".", "conf", "config.yaml"),
			},
			filepath.Join(".", "conf", "config.yaml")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filePath(tt.args.path, tt.args.defaultPath); got != tt.want {
				t.Errorf("filePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
