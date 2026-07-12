package zbank

import (
	"encoding/binary"
	"io"

	"github.com/frogssoldseparately/shippacker/pkg/sreader"
)

type ZBank struct {
	Meta         *Meta
	Drums        *[]*Drum
	Instruments  *[]*Instrument
	SoundEffects *[]*Sfx
	EnvelopeMap  *map[uint32]*Envelope
	SampleMap    *map[uint32]*Sample
}

// Information on how these .zbank files are structured can be found
// in this great writeup by Tharo at https://hackmd.io/6QDQ7l_1T-CExkSsbDCaOA

func NewBankFromStream(f io.Reader, meta *Meta) (*ZBank, error) {
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
	for offset := range sampleMap {
		r.Seek(offset, 0)
		sampleMap[offset] = ReadSample(r)
	}

	return &ZBank{meta, &drums, &instruments, &soundEffects, &envelopeMap, &sampleMap}, nil
}
