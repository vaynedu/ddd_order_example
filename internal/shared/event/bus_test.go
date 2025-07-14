package event_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vaynedu/ddd_order_example/internal/infrastructure/mocks"
	"github.com/vaynedu/ddd_order_example/internal/shared/event"
	"go.uber.org/mock/gomock"
)

// TestEvent 实现Event接口的测试事件
type TestEvent struct {
	name string
}

func (e *TestEvent) Name() string {
	return e.name
}

// TestEventBus_RegisterHandlerAndPublish 测试注册处理器并发布事件
func TestEventBus_RegisterHandlerAndPublish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建mock处理器
	mockHandler := mocks.NewMockHandler(ctrl)
	bus := event.NewEventBus()
	event := &TestEvent{name: "order.created"}

	// 设置预期：处理器应被调用一次
	mockHandler.EXPECT().Handle(gomock.Any(), gomock.Eq(event)).Return(nil).Times(1)

	// 注册处理器并发布事件
	bus.RegisterHandler(event.Name(), mockHandler)
	err := bus.Publish(context.Background(), event)

	// 验证结果
	assert.NoError(t, err)
}

// TestEventBus_PublishWithoutHandlers 测试发布没有处理器的事件
func TestEventBus_PublishWithoutHandlers(t *testing.T) {
	bus := event.NewEventBus()
	event := &TestEvent{name: "order.updated"}

	// 发布没有注册处理器的事件
	err := bus.Publish(context.Background(), event)

	// 验证不会返回错误
	assert.NoError(t, err)
}

// TestEventBus_MultipleHandlers 测试多个处理器接收同一事件
func TestEventBus_MultipleHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建两个mock处理器
	mockHandler1 := mocks.NewMockHandler(ctrl)
	mockHandler2 := mocks.NewMockHandler(ctrl)
	bus := event.NewEventBus()
	event := &TestEvent{name: "payment.succeeded"}

	// 设置每个处理器的预期
	mockHandler1.EXPECT().Handle(gomock.Any(), gomock.Eq(event)).Return(nil).Times(1)
	mockHandler2.EXPECT().Handle(gomock.Any(), gomock.Eq(event)).Return(nil).Times(1)

	// 注册多个处理器
	bus.RegisterHandler(event.Name(), mockHandler1)
	bus.RegisterHandler(event.Name(), mockHandler2)

	// 发布事件
	err := bus.Publish(context.Background(), event)

	// 验证结果
	assert.NoError(t, err)
}
