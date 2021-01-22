package app

import (
	"abs/pkg/logging"
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

// 抛错封装
type Panic struct {
	Err    error
	Stack  []byte
}

// 不要随意改这个, 这个是并发调度方法
// 并发调用服务，每个handler都会传入一个调用逻辑函数
func GoroutineNotPanic(handlers ...func() error) (err error) {
	var (
		wg   sync.WaitGroup
		pErr Panic
	)
	for _, f := range handlers {
		wg.Add(1)
		// 每个函数启动一个协程
		go func(handler func() error) {
			defer func() {
				// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
				// 某个服务调用协程报错，可以在这里打印一些错误日志
				if e := recover(); e != nil {
					err = e.(error)
					pErr = Panic{Err: e.(error), Stack: debug.Stack()}
				}
				wg.Done()
			}()
			// 取第一个报错的handler调用逻辑，并最终向外返回
			e := handler()
			if err == nil && e != nil {
				err = e
			}
			// 解除协程占用Cpu
			return
		}(f)
	}
	// 等待执行完
	wg.Wait()
	// 记录抛出的致命错误
	if pErr.Err != nil {
		msg := fmt.Sprintf("协程组调用致命抛错【Error: %s】 \nStack: %s", pErr.Err, pErr.Stack)
		fmt.Println(msg)
		logging.Error(msg)
	}
	return
}

// 仓库初始化
// 第二种并发调度方法，向上抛出panic错误
func NewPanicGroup() *PanicGroup {
	return &PanicGroup{
		panics: make(chan Panic, 8),
		dones:  make(chan int, 8),
	}
}

type PanicGroup struct {
	panics chan Panic // 协程panic通知信道
	dones  chan int   // 协程完成通知信道
	jobN   int32      // 协程并发数量
}

// go封装调度
func (g *PanicGroup) Go(f func()) *PanicGroup {
	atomic.AddInt32(&g.jobN, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.panics <- Panic{Err: r.(error), Stack: debug.Stack()}
				return
			}
			g.dones <- 1
		}()
		f()
	}()

	// 方便链式调用
	return g
}

// 等待完成
func (g *PanicGroup) Wait(ctx context.Context) error {
	for {
		select {
		case <-g.dones:
			if atomic.AddInt32(&g.jobN, -1) == 0 {
				return nil
			}
		case p := <-g.panics:
			panic(p)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}