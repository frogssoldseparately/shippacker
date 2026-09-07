package zbank

import (
	"encoding/binary"
	"io"

	"github.com/frogssoldseparately/simpleseek/sreader"
)

type ZBank struct {
	Meta         *Meta
	Drums        *[]*Drum
	Instruments  *[]*Instrument
	SoundEffects *[]*Sfx
	EnvelopeMap  *map[uint32]*Envelope
	SampleMap    *map[uint32]*Sample
	LoopMap      *map[uint32]*AdpcmLoop
	BookMap      *map[uint32]*AdpcmBook
}

// Information on how these .zbank files are structured can be found
// in this great writeup by Tharo at https://hackmd.io/6QDQ7l_1T-CExkSsbDCaOA

func NewBankFromStream(f io.Reader, meta *Meta) (*ZBank, error) {
	emptyTunedSample := &TunedSample{0, 0}
	r := sreader.NewSimpleReader(f, binary.BigEndian)
	drumPointerArrayPointer := Read[uint32](r)
	sfxArrayPointer := Read[uint32](r)
	var drums []*Drum
	var instruments []*Instrument
	var soundEffects []*Sfx
	envelopeMap := map[uint32]*Envelope{}
	sampleMap := map[uint32]*Sample{}
	for i := int8(0); i < meta.NumInstruments; i++ {
		r.Seek(0x8+uint32(i)*0x4, 0)
		instHead := r.Seek(Read[uint32](r), 0)
		if instHead != 0 {
			inst := ReadInstrument(r)
			instruments = append(instruments, inst)
			envelopeMap[inst.EnvelopePointer] = nil
			for _, t := range inst.GetTunedSamples() {
				sampleMap[t.SamplePointer] = nil
			}
		} else {
			inst := Instrument{
				ValidByte:              0x1,
				IsRelocated:            0x0,
				NormalRangeLo:          0x0,
				NormalRangeHi:          0x7F,
				AdsrDecayIndex:         0x75,
				EnvelopePointer:        0x0,
				LowPitchTunedSample:    emptyTunedSample,
				NormalPitchTunedSample: emptyTunedSample,
				HighPitchTunedSample:   emptyTunedSample,
			}
			instruments = append(instruments, &inst)
		}
	}
	if drumPointerArrayPointer != 0 {
		for i := int8(0); i < meta.NumDrums; i++ {
			r.Seek(drumPointerArrayPointer+uint32(i)*0x4, 0)
			drumHead := r.Seek(Read[uint32](r), 0)
			if drumHead != 0 {
				drum := ReadDrum(r)
				drums = append(drums, drum)
				envelopeMap[drum.EnvelopePointer] = nil
				sampleMap[drum.TunedSample.SamplePointer] = nil
			} else {
				drum := Drum{}
				drum.NonNull = false
				drums = append(drums, &drum)
			}
		}
	}
	if sfxArrayPointer != 0 {
		r.Seek(sfxArrayPointer, 0)
		for i := int16(0); i < meta.NumSfx; i++ {
			sfx := ReadSfx(r)
			soundEffects = append(soundEffects, sfx)
			sampleMap[sfx.TunedSample.SamplePointer] = nil
		}
	}
	delete(envelopeMap, 0)
	delete(sampleMap, 0)
	for offset := range envelopeMap {
		r.Seek(offset, 0)
		envelopeMap[offset] = ReadEnvelope(r)
	}
	{
		points := []*EnvelopePoint{}
		envelopeMap[0] = &Envelope{Points: &points}
	}
	loopMap := map[uint32]*AdpcmLoop{}
	bookMap := map[uint32]*AdpcmBook{}
	var sampleOffset uint32
	switch meta.SampleBankId1 {
	case 3:
		sampleOffset = 0x4006B0
	case 6:
		sampleOffset = 0x4377E0
	default:
		sampleOffset = 0x0
	}
	for offset := range sampleMap {
		r.Seek(offset, 0)
		sample := ReadSample(r)
		sample.SampleAddress += sampleOffset
		sampleMap[offset] = sample
		r.Seek(sample.LoopPointer, 0)
		loopMap[sample.SampleAddress] = ReadLoop(r)
		r.Seek(sample.BookPointer, 0)
		bookMap[sample.SampleAddress] = ReadBook(r)
	}

	meta.SampleBankId1 = 1

	return &ZBank{meta, &drums, &instruments, &soundEffects, &envelopeMap, &sampleMap, &loopMap, &bookMap}, nil
}
