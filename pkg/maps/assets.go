package maps

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/frogssoldseparately/shippacker/pkg/globals"
)

// A map of sample pointers (relative addresses within then sample banks on ROM) to the asset names used by 2Ship2Harkinian.
type Assets map[uint32]string

// Generate an asset map using mm.o2r and Audio.xml for .zbank -> Soundfont conversion.
func NewAssetMap() (*Assets, error) {
	key := globals.GetAudioXmlKey()
	buf, ok := globals.SampleXmls[key]
	if !ok {
		return nil, fmt.Errorf("Unknown xml key %s", key)
	}
	sampleEntries, err := ParseXML(buf)
	if err != nil {
		return nil, err
	}
	out := Assets{}

	for _, sample := range sampleEntries {
		currentOffset, _ := sample.GetOffset()
		out[currentOffset] = sample.Name
	}
	return &out, nil
}

func SetSample(sampleMap map[string]uint32, filename string) {
	parts := strings.Split(filename, "/")
	outtermost := parts[2]
	sampleMap[outtermost[0:len(outtermost)-len("_META")]] = 0x0
}

func ReadCentralHeader(f *os.File, signature *uint32, compressedSize *uint32, filenameLength *uint16, localHeaderOffset *uint32) (string, error) {
	binary.Read(f, binary.LittleEndian, signature)
	if *signature != 0x02014b50 {
		return "", fmt.Errorf("Unexpected end")
	}
	f.Seek(0x10, 1)
	binary.Read(f, binary.LittleEndian, compressedSize)
	f.Seek(0x4, 1)
	binary.Read(f, binary.LittleEndian, filenameLength)
	f.Seek(0xC, 1)
	binary.Read(f, binary.LittleEndian, localHeaderOffset)
	buf := make([]byte, *filenameLength)
	f.Read(buf)
	return string(buf[:]), nil
}

func ReadUpToFilenameMatch(f *os.File, signature *uint32, compressedSize *uint32, filenameLength *uint16, localHeaderOffset *uint32, match string) (uint32, error) {
	readFiles := uint32(0)
	for {
		filename, err := ReadCentralHeader(f, signature, compressedSize, filenameLength, localHeaderOffset)
		readFiles++
		if err != nil {
			return readFiles, err
		}
		if filename == match {
			return readFiles, nil
		}
	}
}
