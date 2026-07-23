package seq

import (
	"encoding/binary"
	"io"
	"os"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/o2r"
	"github.com/frogssoldseparately/shippacker/pkg/seqcat"
	"github.com/frogssoldseparately/shippacker/pkg/sreader"
	"github.com/frogssoldseparately/shippacker/pkg/swriter"
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

func NewSequenceFromFileWithMeta(seqPath string, metaPath string) (*Sequence, error) {
	if _, err := os.Stat(metaPath); err == nil {
		o2rMeta, err := NewMetaFromFile(metaPath)
		if err != nil {
			return nil, err
		}
		r, err := os.Open(seqPath)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		return NewSequenceFromStream(r, o2rMeta.DisplayName, o2rMeta.BankIds)
	}
	return NewSequenceFromFile(seqPath)
}

func NewSequenceFromFile(path string) (*Sequence, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	moddedName, bankIds, err := ExtractInformationFromPath(path)
	if err != nil {
		return nil, err
	}
	return NewSequenceFromStream(r, moddedName, bankIds)
}

func NewSequenceFromStream(f io.Reader, name string, bankIds *[]byte) (*Sequence, error) {
	nameParts := strings.Split(name, "_")
	songType := nameParts[len(nameParts)-1]
	var cachePolicy uint8 = 0x2
	if seqcat.IsFanfareString(songType) {
		// change to PERSISTENT
		cachePolicy = 0x1
	}
	r := sreader.NewSimpleReader(f, binary.LittleEndian)
	return &Sequence{
		r.GetLength(),
		r.GetBuffer(),
		0x0,
		o2r.SequenceMedium,
		cachePolicy,
		uint32(len(*bankIds)),
		bankIds,
		"custom/music/" + name,
	}, nil
}

func (s *Sequence) GetCompression() uint16 {
	return o2r.Deflate
}

func (s *Sequence) GetFilename() string {
	return s.Path
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
