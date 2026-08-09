package zbank

import (
	"github.com/frogssoldseparately/simpleseek/sreader"
)

type TunedSample struct {
	SamplePointer uint32
	Tuning        float32
}

func ReadTunedSample(r *sreader.SimpleReader) *TunedSample {
	return &TunedSample{
		Read[uint32](r),
		Read[float32](r),
	}
}
