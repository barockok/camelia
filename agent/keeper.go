package agent

import (
	"bytes"
	"io"
	"time"

	"fmt"

	"github.com/barockok/camelia/misc"
)

// Keeper todo
type Keeper struct {
	content     io.ReadWriter
	topic       string
	partition   int64
	offsetStart int64
	offsetEnd   int64
	tsStart     *time.Time
	tsEnd       *time.Time
	counter     int64
}

// Add todo
func (k *Keeper) Add(offset int64, msg misc.JsonAble) error {
	b, err := msg.ToJSON()
	if err != nil {
		return err
	}
	_, err = k.content.Write(b)
	if err != nil {
		return err
	}
	if k.counter == 0 {
		k.offsetStart = offset
	}
	k.counter++
	k.offsetEnd = offset
	return nil
}

// Save Todo
func (k *Keeper) Save(s3instance misc.S3Iface) error {
	return nil
}

func NewKeeper(topic string, partition int64) *Keeper {
	var buf bytes.Buffer
	return &Keeper{topic: topic, partition: partition, content: &buf}
}

func buildObjectPath(k *Keeper) string {
	return fmt.Sprintf("%s/%d/A_%d__Z_%d__O_%d.tar.gz", k.topic, k.partition, k.tsStart.Unix(), k.tsEnd.Unix(), (k.offsetEnd - k.offsetStart))
}
