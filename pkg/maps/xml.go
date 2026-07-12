package maps

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
)

func ParseXML(xmlpath string) ([]*Sample, error) {
	xmlFile, err := os.Open(xmlpath)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()
	buf, err := io.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}
	var root Root
	err = xml.Unmarshal(buf, &root)
	if err != nil {
		return nil, err
	}
	entries := []*Sample{}
	for _, bank := range *root.File.Audio.SamplesContainers {
		offsetAdjustment := uint32(0x0)
		if bank.Bank == "2" {
			offsetAdjustment = 0x538CC0
		}
		for _, sample := range *bank.Samples {
			currentOffset, _ := sample.GetOffset()
			if currentOffset >= 0x3C0F20 {
				sample.Bank = "1"
			} else {
				sample.Bank = bank.Bank
			}
			sample.SetOffset(currentOffset + offsetAdjustment)
			entries = append(entries, &sample)
		}
	}
	return entries, nil
}

type Root struct {
	XMLName xml.Name `xml:"Root"`
	File    *File    `xml:"File"`
}

type File struct {
	XMLName xml.Name `xml:"File"`
	Audio   *Audio   `xml:"Audio"`
}

type Audio struct {
	XMLName           xml.Name            `xml:"Audio"`
	SequenceContainer *SequenceContainer  `xml:"Sequences"`
	SamplesContainers *[]SamplesContainer `xml:"Samples"`
	Soundfonts        *[]Soundfont        `xml:"Soundfont"`
}

type SequenceContainer struct {
	XMLName   xml.Name    `xml:"Sequences"`
	Sequences *[]Sequence `xml:"Sequence"`
}

type Sequence struct {
	XMLName xml.Name `xml:"Sequence"`
}

type SamplesContainer struct {
	XMLName xml.Name  `xml:"Samples"`
	Bank    string    `xml:"Bank,attr"`
	Samples *[]Sample `xml:"Sample"`
}

type Sample struct {
	XMLName xml.Name `xml:"Sample"`
	Name    string   `xml:"Name,attr"`
	Offset  string   `xml:"Offset,attr"`
	Bank    string
}

type Soundfont struct {
	XMLName xml.Name `xml:"Soundfont"`
}

func (s *Sample) GetOffset() (uint32, error) {
	val, err := strconv.ParseUint(s.Offset[2:len(s.Offset)], 16, 32)
	if err != nil {
		return 0x0, err
	}
	return uint32(val), nil
}

func (s *Sample) SetOffset(n uint32) {
	s.Offset = fmt.Sprintf("%#08X", n)
}
