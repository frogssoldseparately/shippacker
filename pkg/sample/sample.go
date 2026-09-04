package sample

import (
	"encoding/binary"
	"io"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/o2r"
	"github.com/frogssoldseparately/shippacker/pkg/zbank"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

type Sample struct {
	RawBinary *[]byte
	Loop      *zbank.AdpcmLoop
	Book      *zbank.AdpcmBook
	Addr      uint32
	Name      string
}

func NewSampleFromStream(fSample io.Reader, addr uint32, name string, am *maps.AssetMap) (*Sample, error) {
	r := sreader.NewSimpleReader(fSample, binary.BigEndian)
	return &Sample{r.GetBuffer(), &zbank.AdpcmLoop{}, &zbank.AdpcmBook{}, addr, name}, nil
}

func (s *Sample) GetCompression() uint16 {
	return o2r.Deflate
}

func (s *Sample) GetFilename() string {
	return "audio/samples/" + s.Name
}

func (s *Sample) GetEndianness() uint32 {
	return o2r.LittleEndian
}

func (s *Sample) GetResourceType() uint32 {
	return o2r.SampleType
}

func (s *Sample) GetResourceVersion() uint32 {
	return o2r.SampleVersion
}

func (s *Sample) GetIsCustom() uint8 {
	return o2r.SampleCustom
}

func (s *Sample) GetBitFlag() uint16 {
	return o2r.DeflateFlags
}

func (s *Sample) GetAttributes() uint32 {
	return o2r.Attributes
}

func (s *Sample) WriteLocalBody(w *swriter.SimpleWriter) error {
	o2r.WriteHeader(s, w)
	Write[uint32](w, 0x2000) // unknown value, sometimes 0x2000, 0x3, 0x1..?
	Write(w, uint32(len(*s.RawBinary)))
	WriteRaw(w, s.RawBinary)
	// Write AdpcmLoop
	Write(w, s.Loop.Start)
	Write(w, s.Loop.End)
	Write(w, s.Loop.Count)
	if s.Loop.Count != 0 {
		Write[uint32](w, 16) // Always 0 or 16
		for _, u16 := range *s.Loop.State {
			Write(w, u16)
		}
	} else {
		Write[uint32](w, 0)
	}
	// Write AdpcmBook
	Write(w, s.Book.Order)
	Write(w, s.Book.NPredictors)
	Write(w, uint32(len(*s.Book.Book))) // independent of 8 * order * npred calculation for some custom samples
	for _, u16 := range *s.Book.Book {
		Write(w, u16)
	}
	return nil
}
