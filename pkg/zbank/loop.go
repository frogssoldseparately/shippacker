package zbank

import (
	"github.com/frogssoldseparately/simpleseek/sreader"
)

type AdpcmLoop struct {
	Start   uint32
	End     uint32
	Count   uint32
	Padding uint32
	State   *[16]uint16
	// State *[32]uint8
}

func ReadLoop(r *sreader.SimpleReader) *AdpcmLoop {
	startV := Read[uint32](r)
	endV := Read[uint32](r)
	countV := Read[uint32](r)
	paddingV := Read[uint32](r)
	stateV := [16]uint16{}
	if countV > 0 {
		for i := range 0x10 {
			stateV[i] = Read[uint16](r)
		}
	}
	return &AdpcmLoop{startV, endV, countV, paddingV, &stateV}
}
