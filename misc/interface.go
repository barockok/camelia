package misc

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// S3Iface todo
type S3Iface interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
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
