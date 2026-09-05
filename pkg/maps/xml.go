package maps

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

func ParseXML(buf []byte) ([]*Sample, error) {
	var root Root
	err := xml.Unmarshal(buf, &root)
	if err != nil {
		return nil, err
	}
	entries := []*Sample{}
	for _, bank := range *root.SamplesContainer {
		offsetAdjustment, err := bank.GetOffset()
		if err != nil {
			return nil, err
		}
		for _, sample := range *bank.Samples {
			sample.Bank = bank.Bank
			currentOffset, _ := sample.GetOffset()
			sample.SetOffset(currentOffset + offsetAdjustment)
			entries = append(entries, &sample)
		}
	}
	return entries, nil
}

type Root struct {
	XMLName          xml.Name            `xml:"Root"`
	SamplesContainer *[]SamplesContainer `xml:"Samples"`
}

type SamplesContainer struct {
	XMLName xml.Name  `xml:"Samples"`
	Bank    string    `xml:"Bank,attr"`
	Offset  string    `xml:"Offset,attr"`
	Samples *[]Sample `xml:"Sample"`
}

func (s *SamplesContainer) GetOffset() (uint32, error) {
	val, err := strconv.ParseUint(s.Offset[2:len(s.Offset)], 16, 32)
	if err != nil {
		return 0x0, err
	}
	return uint32(val), nil
}

type Sample struct {
	XMLName      xml.Name `xml:"Sample"`
	Name         string   `xml:"Name,attr"`
	OriginalName string   `xml:"OriginalName,attr"`
	Offset       string   `xml:"Offset,attr"`
	Bank         string
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
