package task

import (
	"container/list"
	"github.com/tornadoyi/viking/goplus/runtime"
	"math"
	"time"
)

const (
	commandBufferSize			= 1024
)

var (
	monitorOnOff				= true
	tasks						= list.New()
	commands					= make(chan cmd, commandBufferSize)
	checkDelay					= 3 * time.Second
	zombieDuration				= int64(math.MaxInt64)
	zombieCallback				func(*Task)
)

func init() {

	check := func() {
		now := time.Now().UnixNano()

		e := tasks.Front()
		for ; e != nil; {
			n := e.Next()
			t := e.Value.(*Task)
			if t.Terminated() || t.skipMonitor {
				tasks.Remove(e)
			}
			if zombieCallback != nil && now - t.createTime.UnixNano() >= zombieDuration { zombieCallback(t) }
			e = n
		}
	}

	go func() {
		for {
			c := <- commands
			switch c.(type) {
			case *addTaskCmd:  tasks.PushBack(c.(*addTaskCmd).task)
			case *checkCmd:  check()
			}
		}
	}()

	var t *time.Timer
	t = time.AfterFunc(checkDelay, func() {
		commands <- &checkCmd{}
		t.Reset(checkDelay)
	})
}


func NumTask() int { return tasks.Len() }


func SetMonitorOnOff(onOff bool) { monitorOnOff = onOff }

func SetCheckDelay(delay time.Duration) { checkDelay = delay }

func SetZombieDuration(d time.Duration) { zombieDuration = int64(d) }

func SetZombieCallback(f func(*Task)) {
	zombieCallback = func(task *Task) {
		defer runtime.Catch()
		f(task)
	}
}


func onTaskCreate(t *Task) {
	if !monitorOnOff { return }
	commands <- &addTaskCmd{t}
}


type cmd interface {}
type checkCmd struct {}
type addTaskCmd struct {
	task		*Task
}
