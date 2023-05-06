package notifier

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
	
	"github.com/gozelle/logging"
)

var log = logging.Logger("notifier")

const (
	successLvl = iota
	warnLvl
	errorLvl
)

type Notifier interface {
	Infof(format string, a ...any)
	Warnf(format string, a ...any)
	Errorf(format string, a ...any)
	Run()
}

func NewNotify(driver Driver) Notifier {
	return &notifier{
		driver:   driver,
		interval: 6 * time.Second,
	}
}

type message struct {
	level   int
	message string
}

var _ Notifier = (*notifier)(nil)

type notifier struct {
	driver   Driver
	once     sync.Once
	lock     sync.Mutex
	messages []*message
	interval time.Duration
}

func (n *notifier) Run() {
	n.once.Do(func() {
		if n.interval == 0 {
			n.interval = 6 * time.Second
		}
		timer := time.NewTimer(n.interval)
		go func() {
			for {
				select {
				case <-timer.C:
					n.flush()
					timer.Reset(n.interval)
				}
			}
		}()
	})
}

func (n *notifier) flush() {
	n.lock.Lock()
	defer func() {
		n.messages = make([]*message, 0)
		n.lock.Unlock()
	}()
	if len(n.messages) == 0 {
		return
	}
	now := time.Now()
	buf := strings.Builder{}
	var sc, wc, ec int
	_, _ = buf.WriteString(fmt.Sprintf("%s [%s]  \n", now.Format("2006-01-02 15:04:05"), n.interval.String()))
	for _, v := range n.messages {
		emoj := ""
		switch v.level {
		case successLvl:
			emoj = "✅"
			sc++
		case warnLvl:
			emoj = "⚠️"
			wc++
		case errorLvl:
			emoj = "❌"
			sc++
		}
		_, _ = buf.WriteString(fmt.Sprintf("%s %s  \n", emoj, v.message))
	}
	title := ""
	if sc > 0 {
		title += fmt.Sprintf("✅ %d ", sc)
	}
	if wc > 0 {
		title += fmt.Sprintf("⚠️ %d ", wc)
	}
	if ec > 0 {
		title += fmt.Sprintf("❌ %d ", ec)
	}
	title = strings.TrimSpace(title)
	err := n.driver.Push(context.Background(), title, buf.String())
	if err != nil {
		log.Errorf("通知发送失败: %s 内容: %s", err, buf.String())
	}
}

func (n *notifier) push(msg *message) {
	n.lock.Lock()
	defer func() {
		n.lock.Unlock()
	}()
	n.messages = append(n.messages, msg)
}

func (n *notifier) Infof(format string, a ...any) {
	n.push(&message{level: successLvl, message: fmt.Sprintf(format, a...)})
}

func (n *notifier) Warnf(format string, a ...any) {
	n.push(&message{level: warnLvl, message: fmt.Sprintf(format, a...)})
}

func (n *notifier) Errorf(format string, a ...any) {
	n.push(&message{level: errorLvl, message: fmt.Sprintf(format, a...)})
}
