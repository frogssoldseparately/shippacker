package sample

import (
	"encoding/binary"
	"io"

	"github.com/frogssoldseparately/shippacker/pkg/o2r"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

type StubbedSample struct {
	RawBinary *[]byte
	Addr      uint32
	Name      string
}

func ReadShipSample(fSample io.Reader, addr uint32, name string) (*StubbedSample, error) {
	r := sreader.NewSimpleReader(fSample, binary.LittleEndian)
	fullBinary := r.GetBuffer()
	binaryBody := (*fullBinary)[0x40:]
	return &StubbedSample{&binaryBody, addr, name}, nil
}

func (s *StubbedSample) GetCompression() uint16 {
	return o2r.Deflate
}

func (s *StubbedSample) GetFilename() string {
	return "audio/samples/" + s.Name
}

func (s *StubbedSample) GetEndianness() uint32 {
	return o2r.LittleEndian
}

func (s *StubbedSample) GetResourceType() uint32 {
	return o2r.SampleType
}

func (s *StubbedSample) GetResourceVersion() uint32 {
	return o2r.SampleVersion
}

func (s *StubbedSample) GetIsCustom() uint8 {
	return o2r.SampleCustom
}

func (s *StubbedSample) GetBitFlag() uint16 {
	return o2r.DeflateFlags
}

func (s *StubbedSample) GetAttributes() uint32 {
	return o2r.Attributes
}

func (s *StubbedSample) WriteLocalBody(w *swriter.SimpleWriter) error {
	o2r.WriteHeader(s, w)
	WriteRaw(w, s.RawBinary)
	return nil
}
