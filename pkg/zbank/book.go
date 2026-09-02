package zbank

import (
	"github.com/frogssoldseparately/simpleseek/sreader"
)

type AdpcmBook struct {
	Order       uint32
	NPredictors uint32
	// Book        *[]uint8
	Book *[]uint16
}

func ReadBook(r *sreader.SimpleReader) *AdpcmBook {
	orderV := Read[uint32](r)
	npredV := Read[uint32](r)
	bookV := []uint16{}
	lastRead := uint16(0xFFFF)
	lastLastRead := uint16(0xFFFF)
	for lastRead != lastLastRead || lastRead != 0 {
		lastLastRead = lastRead
		lastRead = Read[uint16](r)
		bookV = append(bookV, lastRead)
	}
	bookV = bookV[0 : len(bookV)-2]
	return &AdpcmBook{orderV, npredV, &bookV}
}
