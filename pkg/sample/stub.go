package sample

import (
	"github.com/frogssoldseparately/shippacker/pkg/o2r"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

// TODO: This should eventually should just be removed, and instead parse the sample
// file so we don't have Sample and StubbledSample structures.

type StubbedSample struct {
	RawBinary *[]byte
	Addr      uint32
	Name      string
}

func (s *StubbedSample) GetCompression() uint16 {
	return o2r.Deflate
}

func (s *StubbedSample) GetFilename() string {
	return "audio/samples/" + s.Name
}

func (s *StubbedSample) GetFiletype() string {
	return "Sample"
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
