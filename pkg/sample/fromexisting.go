package sample

import (
	"encoding/binary"
	"io"

	"github.com/frogssoldseparately/simpleseek/sreader"
)

func ReadShipSample(fSample io.Reader, addr uint32, name string) (*StubbedSample, error) {
	r := sreader.NewSimpleReader(fSample, binary.LittleEndian)
	fullBinary := r.GetBuffer()
	binaryBody := (*fullBinary)[0x40:]
	return &StubbedSample{&binaryBody, addr, name}, nil
}
