package observer

// ConfigObserver 是一个配置观察者接口
type ConfigObserver interface {
	UpdateConfig(data string) error
}
