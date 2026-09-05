package soundfont

import (
	"encoding/binary"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/o2r"
	"github.com/frogssoldseparately/shippacker/pkg/zbank"
	"github.com/frogssoldseparately/simpleseek/sreader"
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
	LoopMap      *map[uint32]*zbank.AdpcmLoop
	BookMap      *map[uint32]*zbank.AdpcmBook
	AssetMap     *maps.AssetMap
	Path         string
}

func ReadSoundfont(fSoundfont io.Reader, name string, am *maps.AssetMap, tm *maps.TranslationMap) (*Soundfont, error) {
	bankId, err := getBankFromFontName(name)
	if err != nil {
		return nil, err
	}
	r := sreader.NewSimpleReader(fSoundfont, binary.LittleEndian)
	r.Seek(0x44, 0) // skip the o2r header and bank id
	meta := zbank.ReadBankmeta(r)
	// swap sample banks due to endianness being different
	meta.SampleBankId1 = meta.SampleBankId1 ^ meta.SampleBankId2
	meta.SampleBankId2 = meta.SampleBankId2 ^ meta.SampleBankId1
	meta.SampleBankId1 = meta.SampleBankId1 ^ meta.SampleBankId2
	envMap := map[uint32]*zbank.Envelope{}
	sampleMap := map[uint32]*zbank.Sample{}
	loopMap := map[uint32]*zbank.AdpcmLoop{}
	bookMap := map[uint32]*zbank.AdpcmBook{}

	drums := []*zbank.Drum{}
	instruments := []*zbank.Instrument{}
	soundEffects := []*zbank.Sfx{}
	drumCount := Read[uint32](r)
	instCount := Read[uint32](r)
	sfxCount := Read[uint32](r)
	for range drumCount {
		adsrDecayIndex := Read[uint8](r)
		pan := Read[uint8](r)
		isRelocated := Read[uint8](r)
		// Read Envelope
		envPtr := r.Seek(0, 1)
		envLength := Read[uint32](r)
		var envelope *zbank.Envelope
		if envLength != 0 {
			envelope = zbank.ReadEnvelope(r)
		} else {
			points := []*zbank.EnvelopePoint{}
			envelope = &zbank.Envelope{Points: &points}
		}
		// Read Tuned Sample
		samplePtr := r.Seek(1, 1)
		nameLen := Read[uint32](r)
		assetPath := ReadString(r, nameLen)
		assetName := filepath.Base(assetPath)
		assetAddr, ok := (*tm)[assetName]
		if !ok {
			// assetAddr = 0
			return nil, fmt.Errorf("could not find Ship of Harkinian translation for %s\n", assetName)
		}
		tuning := Read[float32](r)
		// Generate Structures
		envMap[envPtr] = envelope
		sampleMap[samplePtr] = &zbank.Sample{
			BitsAndSize:   0x0,
			SampleAddress: assetAddr,
			LoopPointer:   0x0,
			BookPointer:   0x0,
		}
		tunedSample := zbank.TunedSample{SamplePointer: samplePtr, Tuning: tuning}
		drum := zbank.Drum{
			NonNull:         true,
			AdsrDecayIndex:  adsrDecayIndex,
			Pan:             pan,
			IsRelocated:     isRelocated,
			Unused:          0x0,
			TunedSample:     &tunedSample,
			EnvelopePointer: envPtr,
		}
		drums = append(drums, &drum)
	}
	for range instCount {
		r.Seek(1, 1)
		isRelocated := Read[uint8](r)
		normalRangeLo := Read[uint8](r)
		normalRangeHi := Read[uint8](r)
		adsrDecayIndex := Read[uint8](r)
		// Read Envelope
		envPtr := r.Seek(0, 1)
		envLength := Read[uint32](r)
		var envelope *zbank.Envelope
		if envLength != 0 {
			envelope = zbank.ReadEnvelope(r)
		} else {
			points := []*zbank.EnvelopePoint{}
			envelope = &zbank.Envelope{Points: &points}
		}
		envMap[envPtr] = envelope
		// Read Tuned Samples
		tunedSamples := []*zbank.TunedSample{}
		for range 3 { // Low, Normal, Hi tuned samples
			var tunedSample zbank.TunedSample
			if Read[uint8](r) == 0x1 {
				samplePtr := r.Seek(1, 1)
				nameLen := Read[uint32](r)
				assetPath := ReadString(r, nameLen)
				assetName := filepath.Base(assetPath)
				assetAddr, ok := (*tm)[assetName]
				if !ok {
					// assetAddr = 0
					return nil, fmt.Errorf("could not find Ship of Harkinian translation for %s\n", assetName)
				}
				tuning := Read[float32](r)
				sampleMap[samplePtr] = &zbank.Sample{
					BitsAndSize:   0x0,
					SampleAddress: assetAddr,
					LoopPointer:   0x0,
					BookPointer:   0x0,
				}
				tunedSample = zbank.TunedSample{SamplePointer: samplePtr, Tuning: tuning}
			} else {
				tunedSample = zbank.TunedSample{SamplePointer: 0x0, Tuning: 0x0}
			}
			tunedSamples = append(tunedSamples, &tunedSample)
		}
		inst := zbank.Instrument{
			NonNull:                true,
			IsRelocated:            isRelocated,
			NormalRangeLo:          normalRangeLo,
			NormalRangeHi:          normalRangeHi,
			AdsrDecayIndex:         adsrDecayIndex,
			EnvelopePointer:        envPtr,
			LowPitchTunedSample:    tunedSamples[0],
			NormalPitchTunedSample: tunedSamples[1],
			HighPitchTunedSample:   tunedSamples[2],
		}
		instruments = append(instruments, &inst)
	}
	for range sfxCount {
		var tunedSample zbank.TunedSample
		if Read[uint8](r) == 0x1 {
			samplePtr := r.Seek(1, 1)
			nameLen := Read[uint32](r)
			assetPath := ReadString(r, nameLen)
			assetName := filepath.Base(assetPath)
			assetAddr, ok := (*tm)[assetName]
			if !ok {
				// assetAddr = 0
				return nil, fmt.Errorf("could not find Ship of Harkinian translation for %s\n", assetName)
			}
			tuning := Read[float32](r)
			sampleMap[samplePtr] = &zbank.Sample{
				BitsAndSize:   0x0,
				SampleAddress: assetAddr,
				LoopPointer:   0x0,
				BookPointer:   0x0,
			}
			tunedSample = zbank.TunedSample{SamplePointer: samplePtr, Tuning: tuning}
		} else {
			tunedSample = zbank.TunedSample{SamplePointer: 0x0, Tuning: 0x0}
		}
		sfx := zbank.Sfx{
			TunedSample: &tunedSample,
		}
		soundEffects = append(soundEffects, &sfx)
	}
	return &Soundfont{bankId, meta, &drums, &instruments, &soundEffects, &envMap, &sampleMap, &loopMap, &bookMap, am, name}, nil
}

func NewSoundfontFromBankStreams(fBank io.Reader, fMeta io.Reader, name string, am *maps.AssetMap) (*Soundfont, error) {
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

func NewSoundfontFromBank(bank *zbank.ZBank, name string, am *maps.AssetMap) (*Soundfont, error) {
	bankId, err := getBankFromFontName(name)
	if err != nil {
		return nil, err
	}
	return &Soundfont{bankId, bank.Meta, bank.Drums, bank.Instruments, bank.SoundEffects, bank.EnvelopeMap, bank.SampleMap, bank.LoopMap, bank.BookMap, am, name}, nil
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
		if drum.NonNull {
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
		} else {
			// Write a dummy drum
			Write[uint8](w, 0x0)  // AdsrDecayIndex
			Write[uint8](w, 0x0)  // Pan
			Write[uint8](w, 0x0)  // IsRelocated
			Write[uint32](w, 0x0) // Envelope Length
			Write[uint8](w, 0x1)
			WriteString(w, "audio/samples/Accordion_META", true)
			Write[float32](w, 1)
		}
	}
	return nil
}

func (s *Soundfont) WriteInstruments(w *swriter.SimpleWriter) error {
	for _, inst := range *s.Instruments {
		if inst.NonNull {
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
		} else {
			// Write a dummy instrument
			Write[uint8](w, 0x1)
			Write[uint8](w, 0x0)  // IsRelocated
			Write[uint8](w, 0x0)  // NormalRangeLo
			Write[uint8](w, 0x0)  // NormalRangeHi
			Write[uint8](w, 0x0)  // AdsrDecayIndex
			Write[uint32](w, 0x0) // Envelope length
			Write[uint8](w, 0x0)  // TunedSample
			Write[uint8](w, 0x0)  // TunedSample
			Write[uint8](w, 0x0)  // TunedSample
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
