package misc

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ToI(i interface{}) interface{} { return i }

type KFMessageJSON struct {
	Topic     string
	Partition int
	Offset    int
	Key       []byte
	Value     []byte
}
type KFMessage struct{ kafka.Message }

func (kfm *KFMessage) ToJSON() ([]byte, error) {
	_json := KFMessageJSON{
		Topic:     *kfm.TopicPartition.Topic,
		Partition: ToI(kfm.TopicPartition.Partition).(int),
		Offset:    ToI(kfm.TopicPartition.Offset).(int),
		Key:       kfm.Key,
		Value:     kfm.Value,
	}
	return json.Marshal(_json)
}
