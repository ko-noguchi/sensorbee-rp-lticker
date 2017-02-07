package sensorbee_rp_lticker

import (
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"time"
)

type rpTickerIntervalSource struct {
	ctx      *core.Context
	w        core.Writer
	interval time.Duration
}

func (s *rpTickerIntervalSource) GenerateStream(ctx *core.Context, w core.Writer) error {
	for {
		tuple := core.NewTuple(data.Map{"val": data.Int(1)})
		if err := w.Write(ctx, tuple); err != nil {
			return err
		}
		time.Sleep(s.interval)
	}
	return nil
}

func (s *rpTickerIntervalSource) Stop(ctx *core.Context) error {
	return nil
}

func NewIntervalSource(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Source, error) {
	interval := 1 * time.Second
	if v, ok := params["interval"]; ok {
		i, err := data.ToDuration(v)
		if err != nil {
			return nil, err
		}
		interval = i
	}

	return core.ImplementSourceStop(&rpTickerIntervalSource{
		interval: interval,
	}), nil
}
