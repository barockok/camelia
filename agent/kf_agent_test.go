package agent

import (
	"fmt"
	"testing"
	"time"

	"github.com/barockok/camelia/mocks"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/mock"
)

var (
	c               = &mocks.Customerable{}
	shareMsgHandler = func(msg *kafka.Message) (*kafka.Message, error) { return nil, nil }
	shareErrHandler = func(msg *kafka.Message, err error) {}
	kfa             = NewKFAgent("Tester-agent", []string{"sample-topic"}, c, shareMsgHandler, shareErrHandler)
)

type recorder struct {
	calledTime int
}

func (r *recorder) reset() {
	r.calledTime = 0
}
func (r *recorder) inc() {
	r.calledTime++
}
func (r *recorder) count() int {
	return r.calledTime
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
	msgProcessor := func(msg *kafka.Message) (*kafka.Message, error) { rec.inc(); return msg, nil }

	kfa.MessageHandler = msgProcessor
	kfm := kafka.Message{}
	kfa.processMessage(&kfm)
	if rec.calledTime != 1 {
		t.Errorf("Messaga Handler not triggered when receiving message")
	}

	rec.reset()
	msgProcessor = func(msg *kafka.Message) (*kafka.Message, error) { return nil, fmt.Errorf("Some shit happend") }
	errHandler := func(msg *kafka.Message, err error) { rec.inc() }
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
		t.Errorf("Error handler not triggered when Message handler is panicking")
	}

	rec.reset()

	errHandlerPanic := func(msg *kafka.Message, err error) { panic(fmt.Errorf("Err Handler is Panic")) }
	kfa.MessageErrorHandler = errHandlerPanic
	var panicFromErrHandler error

	func() {
		defer func() {
			if r := recover(); r != nil {
				rec.inc()
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

func TestKFAgent_Start(t *testing.T) {
	rec := recorder{}
	eChan := make(chan kafka.Event)
	defer close(eChan)
	c.On("Events").Return(eChan)
	go kfa.Start()

	msgHandler := func(msg *kafka.Message) (*kafka.Message, error) { rec.inc(); return msg, nil }
	kfa.MessageHandler = msgHandler

	kfm := kafka.Message{}
	eChan <- &kfm
	time.Sleep(200 * time.Millisecond)
	if rec.count() != 1 {
		t.Errorf("Message it not processing when send `kafka.Message` to channel")
	}
	rec.reset()

	tpcs := []kafka.TopicPartition{kafka.TopicPartition{Partition: 1, Offset: 1}}

	c.On("Assign", tpcs).Return(nil)
	assignedPar := kafka.AssignedPartitions{Partitions: tpcs}
	eChan <- assignedPar
	time.Sleep(200 * time.Millisecond)

	c.On("Unassign").Return(nil)
	unAssignedPar := kafka.RevokedPartitions{Partitions: tpcs}
	eChan <- unAssignedPar
	time.Sleep(200 * time.Microsecond)

	c.On("Close").Return(nil)
	eChan <- kafka.Error{}
	time.Sleep(200 * time.Millisecond)
	c.AssertExpectations(t)

}
