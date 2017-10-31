package misc

import (
	"io"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Uploader interface {
	Upload(path string, reader io.Reader) error
}

type JsonAble interface {
	ToJSON() ([]byte, error)
}

type Customerable interface {
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) (err error)
	Events() chan kafka.Event
	Assign(topicPartitions []kafka.TopicPartition) (err error)
	Unassign() (err error)
	Close() (err error)
}

type Configuration interface {
	Configure(map[string]string) error
	Configured() bool
}
