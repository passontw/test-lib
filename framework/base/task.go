package base


type Tasker interface {
	RunTask()
}

type fnEmpty func()

// this is a task without paramters and return
type taskEmpty struct {
	fn fnEmpty
}

func (t *taskEmpty) RunTask() {
	t.fn()
}
