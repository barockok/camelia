package agent

import (
	"fmt"
	"time"

	"github.com/barockok/camelia/misc"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

// KFAgent TODO:
type MessageHandler func(*kafka.Message) (*kafka.Message, error)

// KFAgent TODO:
type MessageErrorHandler func(*kafka.Message, error)

// KFAgent TODO:
type KFAgent struct {
	KFAgentID           string
	Topics              []string
	consumer            misc.Customerable
	logger              *log.Entry
	MessageHandler      MessageHandler
	MessageErrorHandler MessageErrorHandler
}

// KFAgent TODO:
func NewKFAgent(agentID string, topics []string, c misc.Customerable, msgHandler MessageHandler, errHandler MessageErrorHandler) *KFAgent {
	a := KFAgent{
		MessageHandler:      msgHandler,
		MessageErrorHandler: errHandler,
		Topics:              topics,
		consumer:            c,
		logger:              log.WithFields(log.Fields{"agent_id": agentID, "aspect": "kf-agent"}),
	}
	return &a
}

// KFAgent TODO:
func (a *KFAgent) Init() error {
	err := a.consumer.SubscribeTopics(a.Topics, a.rebalaceCB)
	if err != nil {
		return err
	}
	return nil
}

// KFAgent TODO:
func (a *KFAgent) Start() {
	defer func() {
		if r := recover(); r != nil {
			a.logger.Infof("Recovering from Consumer crash, %v", r)
		}
	}()

	for ev := range a.consumer.Events() {
		switch e := ev.(type) {
		case kafka.AssignedPartitions:
			a.logger.Infof("Receiving assigned partition %v", e)
			a.consumer.Assign(e.Partitions)
		case kafka.RevokedPartitions:
			a.logger.Infof("Revoked Partition")
			a.consumer.Unassign()
		case *kafka.Message:
			t := e.TopicPartition
			lf := log.Fields{"topic": t.Topic, "offset": t.Offset, "partition": t.Partition, "partition_key": string(e.Key), "value": string(e.Value)}
			log.WithFields(lf).Infof("Recieving Message")
			a.processMessage(e)
		case kafka.PartitionEOF:
			a.logger.Infof("Reached EOF, %v", e)
		case kafka.Error:
			fields := log.Fields{"error": "1", "error_code": fmt.Sprintf("%d", e.Code())}
			a.logger.WithFields(fields).Error(e.String())
			a.Stop()
		default:
			a.logger.Infof("Receiving Default handler, %v", e)
		}
	}
}

// KFAgent TODO:
func (a *KFAgent) Stop() {
	a.consumer.Close()
}

func (a *KFAgent) rebalaceCB(c *kafka.Consumer, e kafka.Event) (err error) {
	a.logger.Info("rebalancing...")
	return nil
}

func (a *KFAgent) processMessage(e *kafka.Message) {
	t := e.TopicPartition
	lf := log.Fields{"topic": t.Topic, "offset": t.Offset, "partition": t.Partition, "partition_key": e.Key}

	// Topest defer, to log that message have been proceed.
	// and top of that, this one is to recover when MessageErrorHandler is panic.
	// it'll log with panic flag, and continue the panic so the initiator (method Start).
	defer func() {
		log.WithFields(lf).WithField("ts", time.Now().UTC().UnixNano()).Infof("Processing done")
		if r := recover(); r != nil {
			a.logger.WithFields(lf).WithError(r.(error)).Errorf("Receiving panic")
			panic(r.(error))
		}
	}()

	log.WithFields(lf).WithField("ts", time.Now().UTC().UnixNano()).Infof("Processing message")

	// this is closure to wrap Message handler,
	// if message handler is panicking, then it'll recover and forward the message to MessageErrorHandler
	wrapHandler := func() error {
		defer func() {
			if r := recover(); r != nil {
				log.WithFields(lf).WithError(r.(error)).WithField("panic", 1).Error("Could not process message")
				a.MessageErrorHandler(e, fmt.Errorf("%v", r))
				log.WithFields(lf).Info("Message forwarded to errorHandler")
			}
		}()
		_, err := a.MessageHandler(e)
		return err
	}
	if err := wrapHandler(); err != nil {
		log.WithFields(lf).WithError(err).Error("Could not process message")
		a.MessageErrorHandler(e, err)
		log.WithFields(lf).Info("Message forwarded to errorHandler")
	}
}
