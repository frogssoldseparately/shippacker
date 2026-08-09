package zbank

import (
	"github.com/frogssoldseparately/simpleseek/sreader"
)

type Sfx struct {
	TunedSample *TunedSample
}

func ReadSfx(r *sreader.SimpleReader) *Sfx {
	return &Sfx{
		ReadTunedSample(r),
	}
}
