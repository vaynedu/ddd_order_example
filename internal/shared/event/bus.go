// internal/shared/event/bus.go
package event

import (
	"context"
	"sync"
)

// 事件接口
type Event interface {
	Name() string
}

// 事件处理器接口
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// 事件总线
type EventBus struct {
	handlers map[string][]Handler
	mutex    sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
	}
}

// 注册事件处理器
func (b *EventBus) RegisterHandler(eventName string, handler Handler) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

// 发布事件
func (b *EventBus) Publish(ctx context.Context, event Event) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	eventName := event.Name()
	if handlers, ok := b.handlers[eventName]; ok {
		// 异步处理事件
		var wg sync.WaitGroup
		for _, handler := range handlers {
			wg.Add(1)
			go func(h Handler) {
				defer wg.Done()
				_ = h.Handle(ctx, event) // 实际生产中应记录错误
			}(handler)
		}
		wg.Wait()
	}

	return nil
}
