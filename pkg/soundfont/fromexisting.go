package soundfont

import (
	"encoding/binary"
	"fmt"
	"io"
	"path/filepath"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/sample"
	"github.com/frogssoldseparately/shippacker/pkg/zbank"
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func ReadSoundfont(fSoundfont io.Reader, name string, zipWriter *swriter.SimpleZipWriter, gsm *maps.GameSampleMap) (*Soundfont, error) {
	bankId, err := getBankFromFontName(name)
	if err != nil {
		return nil, err
	}
	r := sreader.NewSimpleReader(fSoundfont, binary.LittleEndian)
	r.Seek(0x44, 0) // skip the o2r header and bank id
	meta := zbank.ReadBankmeta(r)
	// Swap values, because meta reads big endian, when soundfonts are little endian
	meta.SampleBankId1 ^= meta.SampleBankId2
	meta.SampleBankId2 ^= meta.SampleBankId1
	meta.SampleBankId1 ^= meta.SampleBankId2
	meta.NumInstruments ^= meta.NumDrums
	meta.NumDrums ^= meta.NumInstruments
	meta.NumInstruments ^= meta.NumDrums
	meta.NumSfx = meta.NumSfx>>8 | meta.NumSfx&0xff<<8
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
		tuning := Read[float32](r)
		nonNull := true
		if len(assetPath) > 0 {
			assetName := filepath.Base(assetPath)
			assetAddr, ok := (*gsm.ByName)[assetName]
			if !ok {
				if assetAddr, err = sample.InjectSampleByName(zipWriter, assetName, gsm); err != nil {
					return nil, fmt.Errorf("2ship2harkinian translation for Ship of Harkinian sample \"%s\" could not be found\n", assetName)
				}
			}
			sampleMap[samplePtr] = &zbank.Sample{
				BitsAndSize:   0x0,
				SampleAddress: assetAddr,
				LoopPointer:   0x0,
				BookPointer:   0x0,
			}
		} else {
			nonNull = false
			samplePtr = 0x0
		}
		tunedSample := zbank.TunedSample{SamplePointer: samplePtr, Tuning: tuning}
		// Generate Structures
		envMap[envPtr] = envelope
		drum := zbank.Drum{
			NonNull:         nonNull,
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
		validByte := Read[uint8](r)
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
				assetAddr, ok := (*gsm.ByName)[assetName]
				if !ok {
					if assetAddr, err = sample.InjectSampleByName(zipWriter, assetName, gsm); err != nil {
						return nil, fmt.Errorf("2ship2harkinian translation for Ship of Harkinian sample \"%s\" could not be found\n", assetName)
					}
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
			ValidByte:              validByte,
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
			assetAddr, ok := (*gsm.ByName)[assetName]
			if !ok {
				if assetAddr, err = sample.InjectSampleByName(zipWriter, assetName, gsm); err != nil {
					return nil, fmt.Errorf("2ship2harkinian translation for Ship of Harkinian sample \"%s\" could not be found\n", assetName)
				}
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
	return &Soundfont{
		BankId:        bankId,
		Meta:          meta,
		Drums:         &drums,
		Instruments:   &instruments,
		SoundEffects:  &soundEffects,
		EnvelopeMap:   &envMap,
		SampleMap:     &sampleMap,
		LoopMap:       &loopMap,
		BookMap:       &bookMap,
		GameSampleMap: gsm,
		Path:          name,
	}, nil
}
