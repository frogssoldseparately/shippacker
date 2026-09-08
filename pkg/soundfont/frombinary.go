package soundfont

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/maps"
	"github.com/frogssoldseparately/shippacker/pkg/zbank"
)

func NewSoundfontFromBankStreams(fBank io.Reader, fMeta io.Reader, name string, gsm *maps.GameSampleMap) (*Soundfont, error) {
	meta, err := zbank.NewBankmetaFromStream(fMeta)
	if err != nil {
		return nil, err
	}
	bank, err := zbank.NewBankFromStream(fBank, meta)
	if err != nil {
		return nil, err
	}
	return NewSoundfontFromBank(bank, name, gsm)
}

func NewSoundfontFromBank(bank *zbank.ZBank, name string, gsm *maps.GameSampleMap) (*Soundfont, error) {
	bankId, err := getBankFromFontName(name)
	if err != nil {
		return nil, err
	}
	return &Soundfont{
		BankId:        bankId,
		Meta:          bank.Meta,
		Drums:         bank.Drums,
		Instruments:   bank.Instruments,
		SoundEffects:  bank.SoundEffects,
		EnvelopeMap:   bank.EnvelopeMap,
		SampleMap:     bank.SampleMap,
		LoopMap:       bank.LoopMap,
		BookMap:       bank.BookMap,
		GameSampleMap: gsm,
		Path:          name,
	}, nil
}

func getBankFromFontName(name string) (uint32, error) {
	nameParts := strings.Split(name, "_")
	v, err := strconv.ParseUint(nameParts[len(nameParts)-1], 10, 32)
	if err != nil {
		return 0x0, fmt.Errorf("its bank name could not be parsed")
	}
	return uint32(v), nil
}
