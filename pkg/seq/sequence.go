package seq

import (
	"github.com/frogssoldseparately/shippacker/pkg/o2r"

	"github.com/frogssoldseparately/simpleseek/swriter"
)

type Sequence struct {
	Size        uint32
	RawBinary   *[]byte
	SequenceNum uint8
	Medium      uint8
	CachePolicy uint8
	NumFonts    uint32
	FontIndices *[]byte
	Path        string
}

func (s *Sequence) GetCompression() uint16 {
	return o2r.Deflate
}

func (s *Sequence) GetFilename() string {
	return s.Path
}

func (s *Sequence) GetFiletype() string {
	return "Sequence"
}

func (s *Sequence) GetEndianness() uint32 {
	return o2r.LittleEndian
}

func (s *Sequence) GetResourceType() uint32 {
	return o2r.SequenceType
}

func (s *Sequence) GetResourceVersion() uint32 {
	return o2r.SequenceVersion
}

func (s *Sequence) GetIsCustom() uint8 {
	return o2r.SequenceCustom
}

func (s *Sequence) GetBitFlag() uint16 {
	return o2r.DeflateFlags
}

func (s *Sequence) GetAttributes() uint32 {
	return o2r.Attributes
}

func (s *Sequence) WriteLocalBody(w *swriter.SimpleWriter) error {
	o2r.WriteHeader(s, w)
	Write(w, s.Size)
	WriteRaw(w, s.RawBinary)
	Write(w, s.SequenceNum)
	Write(w, s.Medium)
	Write(w, s.CachePolicy)
	Write(w, s.NumFonts)
	WriteRaw(w, s.FontIndices)
	return nil
}
