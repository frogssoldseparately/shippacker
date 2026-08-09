package zbank

import (
	"github.com/frogssoldseparately/simpleseek/sreader"
)

type Envelope struct {
	Points *[]*EnvelopePoint
}

type EnvelopePoint struct {
	Delay int16
	Arg   int16
}

func ReadEnvelopePoint(r *sreader.SimpleReader) *EnvelopePoint {
	return &EnvelopePoint{
		Read[int16](r),
		Read[int16](r),
	}
}

func ReadEnvelope(r *sreader.SimpleReader) *Envelope {
	var points []*EnvelopePoint
	for {
		point := ReadEnvelopePoint(r)
		points = append(points, point)
		if point.Delay == -1 {
			break
		}
	}
	return &Envelope{
		&points,
	}
}
