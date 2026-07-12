package zbank

import (
	"encoding/binary"
	"io"

	"github.com/frogssoldseparately/shippacker/pkg/sreader"
)

type Meta struct {
	Medium         int8
	CachePolicy    int8
	SampleBankId1  int8
	SampleBankId2  int8
	NumInstruments int8
	NumDrums       int8
	NumSfx         int16
}

func NewBankmetaFromStream(f io.Reader) (*Meta, error) {
	r := sreader.NewSimpleReader(f, binary.BigEndian)
	return &Meta{
		Read[int8](r),
		Read[int8](r),
		Read[int8](r),
		Read[int8](r),
		Read[int8](r),
		Read[int8](r),
		Read[int16](r),
	}, nil
}
