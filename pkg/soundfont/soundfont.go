package soundfont

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/o2r"
	"github.com/frogssoldseparately/shippacker/pkg/zbank"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

type Soundfont struct {
	BankId       uint32
	Meta         *zbank.Meta
	Drums        *[]*zbank.Drum
	Instruments  *[]*zbank.Instrument
	SoundEffects *[]*zbank.Sfx
	EnvelopeMap  *map[uint32]*zbank.Envelope
	SampleMap    *map[uint32]*zbank.Sample
	AssetMap     *maps.Assets
	Path         string
}

func NewSoundfontFromBankStreams(fBank io.Reader, fMeta io.Reader, name string, am *maps.Assets) (*Soundfont, error) {
	meta, err := zbank.NewBankmetaFromStream(fMeta)
	if err != nil {
		return nil, err
	}
	bank, err := zbank.NewBankFromStream(fBank, meta)
	if err != nil {
		return nil, err
	}
	return NewSoundfontFromBank(bank, name, am)
}

func NewSoundfontFromBank(bank *zbank.ZBank, name string, am *maps.Assets) (*Soundfont, error) {
	bankId, err := getBankFromFontName(name)
	if err != nil {
		return nil, err
	}
	return &Soundfont{bankId, bank.Meta, bank.Drums, bank.Instruments, bank.SoundEffects, bank.EnvelopeMap, bank.SampleMap, am, "custom/fonts/" + name}, nil
}

func (s *Soundfont) GetCompression() uint16 {
	return o2r.Deflate
}

func (s *Soundfont) GetFilename() string {
	return s.Path
}

func (s *Soundfont) GetEndianness() uint32 {
	return o2r.LittleEndian
}

func (s *Soundfont) GetResourceType() uint32 {
	return o2r.SoundfontType
}

func (s *Soundfont) GetResourceVersion() uint32 {
	return o2r.SoundfontVersion
}

func (s *Soundfont) GetIsCustom() uint8 {
	return o2r.SoundfontCustom
}

func (s *Soundfont) GetBitFlag() uint16 {
	return o2r.DeflateFlags
}

func (s *Soundfont) GetAttributes() uint32 {
	return o2r.Attributes
}

func (s *Soundfont) WriteLocalBody(w *swriter.SimpleWriter) error {
	o2r.WriteHeader(s, w)
	s.WriteBankInfo(w)
	if err := s.WriteDrums(w); err != nil {
		return err
	}
	if err := s.WriteInstruments(w); err != nil {
		return err
	}
	return s.WriteSfx(w)
}

func (s *Soundfont) WriteBankInfo(w *swriter.SimpleWriter) {
	Write(w, s.BankId)
	Write(w, s.Meta.Medium)
	Write(w, s.Meta.CachePolicy)
	Write(w, s.Meta.SampleBankId2)
	Write(w, s.Meta.SampleBankId1)
	Write(w, s.Meta.NumDrums)
	Write(w, s.Meta.NumInstruments)
	Write(w, s.Meta.NumSfx)
	Write(w, uint32(len(*s.Drums)))
	Write(w, uint32(len(*s.Instruments)))
	Write(w, uint32(len(*s.SoundEffects)))
}

func (s *Soundfont) WriteEnvelopeEntry(w *swriter.SimpleWriter, ptr uint32) error {
	if e, ok := (*s.EnvelopeMap)[ptr]; ok {
		Write(w, uint32(len(*e.Points)))
		for _, p := range *e.Points {
			Write(w, p.Delay)
			Write(w, p.Arg)
		}
	} else {
		return fmt.Errorf("invalid envelope pointer of %08X\n", ptr)
	}
	return nil
}

func (s *Soundfont) WriteTunedSample(w *swriter.SimpleWriter, ts *zbank.TunedSample) error {
	sample := (*s.SampleMap)[ts.SamplePointer]
	assetName, ok := (*s.AssetMap)[sample.SampleAddress]
	if !ok {
		return fmt.Errorf("invalid sample address of %08X\n", sample.SampleAddress)
	}
	WriteString(w, "audio/samples/"+assetName, true)
	Write(w, ts.Tuning)
	return nil
}

func (s *Soundfont) WriteDrums(w *swriter.SimpleWriter) error {
	for _, drum := range *s.Drums {
		Write(w, drum.AdsrDecayIndex)
		Write(w, drum.Pan)
		Write(w, drum.IsRelocated)
		if err := s.WriteEnvelopeEntry(w, drum.EnvelopePointer); err != nil {
			return err
		}
		Write[uint8](w, 0x1)
		if err := s.WriteTunedSample(w, drum.TunedSample); err != nil {
			return err
		}
	}
	return nil
}

func (s *Soundfont) WriteInstruments(w *swriter.SimpleWriter) error {
	for _, inst := range *s.Instruments {
		Write[uint8](w, 0x1)
		Write(w, inst.IsRelocated)
		Write(w, inst.NormalRangeLo)
		Write(w, inst.NormalRangeHi)
		Write(w, inst.AdsrDecayIndex)
		if err := s.WriteEnvelopeEntry(w, inst.EnvelopePointer); err != nil {
			return err
		}
		for _, tunedSample := range inst.GetTunedSamples() {
			if tunedSample.SamplePointer != 0x0 {
				Write[uint8](w, 0x1)
				Write[uint8](w, 0x1)
				if err := s.WriteTunedSample(w, tunedSample); err != nil {
					return err
				}
			} else {
				Write[uint8](w, 0x0)
			}
		}
	}
	return nil
}

func (s *Soundfont) WriteSfx(w *swriter.SimpleWriter) error {
	for _, sfx := range *s.SoundEffects {
		tunedSample := sfx.TunedSample
		if tunedSample.SamplePointer != 0x0 {
			Write[uint8](w, 0x1)
			Write[uint8](w, 0x1)
			if err := s.WriteTunedSample(w, tunedSample); err != nil {
				return err
			}
		} else {
			Write[uint8](w, 0x0)
		}
	}
	return nil
}

func getBankFromFontName(name string) (uint32, error) {
	nameParts := strings.Split(name, "_")
	v, err := strconv.ParseUint(nameParts[len(nameParts)-1], 10, 32)
	if err != nil {
		return 0x0, fmt.Errorf("its bank name could not be parsed")
	}
	return uint32(v), nil
}
