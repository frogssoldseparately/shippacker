package sample

import (
	"encoding/binary"
	"io"

	"github.com/frogssoldseparately/shippacker/pkg/zbank"
	"github.com/frogssoldseparately/simpleseek/sreader"
)

func NewSampleFromStream(fSample io.Reader, addr uint32, name string) (*Sample, error) {
	r := sreader.NewSimpleReader(fSample, binary.BigEndian)
	return &Sample{r.GetBuffer(), &zbank.AdpcmLoop{}, &zbank.AdpcmBook{}, addr, name}, nil
}
