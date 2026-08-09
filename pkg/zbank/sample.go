package zbank

import (
	"github.com/frogssoldseparately/simpleseek/sreader"
)

type Sample struct {
	BitsAndSize   uint32
	SampleAddress uint32
	LoopPointer   uint32
	BookPointer   uint32
}

func ReadSample(r *sreader.SimpleReader) *Sample {
	return &Sample{
		Read[uint32](r),
		Read[uint32](r),
		Read[uint32](r),
		Read[uint32](r),
	}
}
