package zbank

import (
	"github.com/frogssoldseparately/shippacker/pkg/sreader"
)

type Sfx struct {
	TunedSample *TunedSample
}

func ReadSfx(r *sreader.SimpleReader) *Sfx {
	return &Sfx{
		ReadTunedSample(r),
	}
}
