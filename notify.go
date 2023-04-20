package notifier

import (
	"sync"
	"time"
)

type Notify interface {
	Success(format string, a ...any)
	Warn(format string, a ...any)
	Error(format string, a ...any)
	Run()
}

func NewNotify() Notify {
	return &notify{}
}

var _ Notify = (*notify)(nil)

type notify struct {
	driver Driver
	timer  time.Timer
	once   sync.Once
}

func (n *notify) Run() {
	n.once.Do(func() {
	
	})
}

func (n *notify) Success(format string, a ...any) {
	//TODO implement me
	panic("implement me")
}

func (n *notify) Warn(format string, a ...any) {
	//TODO implement me
	panic("implement me")
}

func (n *notify) Error(format string, a ...any) {
	//TODO implement me
	panic("implement me")
}
