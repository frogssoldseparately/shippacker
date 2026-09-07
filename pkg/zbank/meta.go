package zbank

import (
	"encoding/binary"
	"io"

	"github.com/frogssoldseparately/simpleseek/sreader"
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
	return ReadBankmeta(r), nil
}

func ReadBankmeta(r *sreader.SimpleReader) *Meta {
	medium := Read[int8](r)
	cachePolicy := Read[int8](r)
	sampleBankId1 := Read[int8](r)
	sampleBankId2 := Read[int8](r)
	numInstruments := Read[int8](r)
	numDrums := Read[int8](r)
	var numSfx int16
	if r.Seek(0, 1) != r.GetLength() {
		numSfx = Read[int16](r)
	} else {
		numSfx = 0x0
	}
	return &Meta{
		medium,
		cachePolicy,
		sampleBankId1,
		sampleBankId2,
		numInstruments,
		numDrums,
		numSfx,
	}
}
