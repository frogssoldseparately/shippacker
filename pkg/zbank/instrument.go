package zbank

import (
	"github.com/frogssoldseparately/simpleseek/sreader"
)

type Instrument struct {
	IsRelocated            uint8
	NormalRangeLo          uint8
	NormalRangeHi          uint8
	AdsrDecayIndex         uint8
	EnvelopePointer        uint32
	LowPitchTunedSample    *TunedSample
	NormalPitchTunedSample *TunedSample
	HighPitchTunedSample   *TunedSample
}

func (i *Instrument) GetTunedSamples() []*TunedSample {
	return []*TunedSample{i.LowPitchTunedSample, i.NormalPitchTunedSample, i.HighPitchTunedSample}
}

func ReadInstrument(r *sreader.SimpleReader) *Instrument {
	return &Instrument{
		Read[uint8](r),
		Read[uint8](r),
		Read[uint8](r),
		Read[uint8](r),
		Read[uint32](r),
		ReadTunedSample(r),
		ReadTunedSample(r),
		ReadTunedSample(r),
	}
}
