package sleep

import (
	"context"
	"time"
)

type Sleep struct {
	context context.Context
}

func New(context context.Context) *Sleep {
	s := &Sleep{
		context: context,
	}
	return s
}

func (s *Sleep) IsRunning() bool {
	return s.context.Err() == nil
}

func (s *Sleep) Sleep(seconds int) {
	d := seconds * 2
	for i := 0; i < d && s.IsRunning(); i++ {
		time.Sleep(time.Second / 2)
	}
}
