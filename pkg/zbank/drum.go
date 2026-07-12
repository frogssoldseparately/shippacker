package zbank

import (
	"github.com/frogssoldseparately/shippacker/pkg/sreader"
)

type Drum struct {
	AdsrDecayIndex  uint8
	Pan             uint8
	IsRelocated     uint8
	Unused          uint8
	TunedSample     *TunedSample
	EnvelopePointer uint32
}

func ReadDrum(r *sreader.SimpleReader) *Drum {
	return &Drum{
		Read[uint8](r),
		Read[uint8](r),
		Read[uint8](r),
		Read[uint8](r),
		ReadTunedSample(r),
		Read[uint32](r),
	}
}
