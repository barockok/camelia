package agent

import (
	"testing"

	"fmt"

	"github.com/barockok/camelia/agent/mocks"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/mock"
)

var (
	c          = &mocks.Customerable{}
	msgHandler = func(msg *kafka.Message) (*kafka.Message, error) { return nil, nil }
	errHandler = func(msg *kafka.Message, err error) {}
	kfa        = NewKFAgent("Tester-agent", []string{"sample-topic"}, c, msgHandler, errHandler)
)

type recorder struct {
	calledTime int
}

func (r *recorder) reset() {
	r.calledTime = 0
}
func TestKFAgent_Init(t *testing.T) {
	c.On("SubscribeTopics", []string{"sample-topic"}, mock.Anything).Return(nil)
	kfa.Init()
	c.AssertExpectations(t)
}

func TestKFAgent_Stop(t *testing.T) {
	c.On("Close").Return(nil)
	kfa.Stop()
	c.AssertExpectations(t)
}

func TestKFAgent_processMessage(t *testing.T) {
	rec := recorder{}
	msgProcessor := func(msg *kafka.Message) (*kafka.Message, error) { rec.calledTime++; return msg, nil }

	kfa.MessageHandler = msgProcessor
	kfm := kafka.Message{}
	kfa.processMessage(&kfm)
	if rec.calledTime != 1 {
		t.Errorf("Messaga Handler not triggered when receiving message")
	}

	rec.reset()
	msgProcessor = func(msg *kafka.Message) (*kafka.Message, error) { return nil, fmt.Errorf("Some shit happend") }
	errHandler := func(msg *kafka.Message, err error) { rec.calledTime++ }
	kfa.MessageHandler = msgProcessor
	kfa.MessageErrorHandler = errHandler

	kfa.processMessage(&kfm)
	if rec.calledTime != 1 {
		t.Errorf("Error handler not triggerd when Message handler fail")
	}

	rec.reset()
	msgProcessor = func(msg *kafka.Message) (*kafka.Message, error) { panic(fmt.Errorf("Some shit happend")) }
	kfa.MessageHandler = msgProcessor
	func() {
		defer recover()
		kfa.processMessage(&kfm)
	}()

	if rec.calledTime != 1 {
		t.Errorf("Error handler not triggerd when Message handler panicking")
	}

	rec.reset()

	errHandlerPanic := func(msg *kafka.Message, err error) { panic(fmt.Errorf("Err Handler is Panic")) }
	kfa.MessageErrorHandler = errHandlerPanic
	var panicFromErrHandler error

	func() {
		defer func() {
			if r := recover(); r != nil {
				rec.calledTime++
				panicFromErrHandler = r.(error)
			}
		}()
		kfa.processMessage(&kfm)
	}()

	gotCallExpect := rec.calledTime
	callExpect := 1

	if gotCallExpect != callExpect {
		t.Fatalf("Panic never raised ; expected : %v & got : %v", callExpect, gotCallExpect)
	}

	gotPanicMessage := panicFromErrHandler.Error()
	panicMessage := "Err Handler is Panic"

	if gotPanicMessage != panicMessage {
		t.Errorf("Panic raised, but the message is unexpected, expect : %v & got %v", panicMessage, gotPanicMessage)
	}

}
