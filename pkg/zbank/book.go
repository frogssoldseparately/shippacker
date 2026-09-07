package zbank

import (
	"github.com/frogssoldseparately/simpleseek/sreader"
)

type AdpcmBook struct {
	Order       uint32
	NPredictors uint32
	Book        *[]uint16
}

func ReadBook(r *sreader.SimpleReader) *AdpcmBook {
	orderV := Read[uint32](r)
	npredV := Read[uint32](r)
	bookV := []uint16{}
	consecutiveZeroCount := 0
	nonZeroSeen := false
	for !nonZeroSeen || (len(bookV)-4)%8 != 0 || consecutiveZeroCount != 4 {
		v := Read[uint16](r)
		bookV = append(bookV, v)
		if v == 0 {
			consecutiveZeroCount++
		} else {
			consecutiveZeroCount = 0
			nonZeroSeen = true
		}
	}
	bookV = bookV[0 : len(bookV)-4]
	return &AdpcmBook{orderV, npredV, &bookV}
}
