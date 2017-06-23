package agent

import (
	"fmt"
	"testing"
	"time"

	"github.com/barockok/camelia/mocks"
)

func Test_buildObjectPath(t *testing.T) {
	type args struct {
		k *Keeper
	}
	tsEnd := time.Now()
	tsStart := tsEnd.Add(-2 * time.Second)

	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Simple",
			args: args{k: &Keeper{
				topic: "nice-topic", partition: 1,
				offsetStart: 250, offsetEnd: 788,
				tsEnd:   &tsEnd,
				tsStart: &tsStart}},
			want: fmt.Sprintf("nice-topic/1/A_%d__Z_%d__O_%d.tar.gz", tsStart.Unix(), tsEnd.Unix(), 538),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildObjectPath(tt.args.k); got != tt.want {
				t.Errorf("buildObjectPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_Add(t *testing.T) {
	jsonByte := []byte("this is sample of json byte")
	buffer := &mocks.ReaderWriter{}
	k := &Keeper{content: buffer}
	kfm := mocks.JsonAble{}
	kfm.On("ToJSON").Return(jsonByte, nil)
	buffer.On("Write", jsonByte).Return(1, nil)

	if err := k.Add(100, &kfm); err != nil {
		t.Errorf("can't add message to buffer %v", err)
	}

	if !c.AssertExpectations(t) {
		t.Errorf("Dependencies operation not satisfied")
	}

	if got := k.offsetStart; got != 100 {
		t.Errorf("Add not updated the offsetStart ; expected %v got %v", 100, got)
	}

	k = &Keeper{content: buffer}

	if err := k.Add(0, &kfm); err != nil {
		t.Errorf("can't add message to buffer %v", err)
	}

	if got := k.offsetStart; got != 0 {
		t.Errorf("Add with index 0, not set startOffset with 0; expected %v got %v", 0, got)
	}
	if got := k.offsetEnd; got != 0 {
		t.Errorf("Add with index 0, not set endOffset with 0; expected %v got %v", 0, got)
	}

	if err := k.Add(1, &kfm); err != nil {
		t.Errorf("can't add message to buffer %v", err)
	}

	if k.offsetStart != 0 {
		t.Errorf("Add more byte, keep update the offsetStart; e : %v, g: %v", 0, k.offsetStart)
	}

	if got := k.offsetEnd; got != 1 {
		t.Errorf("Add more byte, not updated the offsetEnd; e : %v, g: %v", 1, got)
	}
}
