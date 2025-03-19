package gamelogic

func NewEventExecuteMySql(traceId string) *EventExecuteMySql {
	return &EventExecuteMySql{
		traceId: traceId,
	}
}

var _ IEvent = (*EventExecuteMySql)(nil)

// EventExecuteMySql 数据更新或者插入事件
type EventExecuteMySql struct {
	traceId string
}

// HandleEvent todo:暂时不写入数据库 待使用时再完善
func (e *EventExecuteMySql) HandleEvent() {
}
