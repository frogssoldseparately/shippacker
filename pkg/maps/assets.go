package maps

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

// A map of sample pointers (relative addresses within then sample banks on ROM) to the asset names used by 2Ship2Harkinian.
type Assets map[uint32]string

// Tracks points of interest within the mm.o2r file.
type ReadPoints struct {
	SoundfontsEnd                 uint32 // the end of audio/fonts/Soundfont_9
	SamplesBegin                  uint32 // the start of audio/samples/Accordion_META
	SamplesEnd                    uint32 // the end of audio/samples/sample_2_005412F0_META
	SoundfontsCentralDirectoryEnd uint32
	SamplesCentralDirectoryBegin  uint32
	CentralDirectoryEnd           uint32
	SkippedLocalSize              uint32 // Size in bytes of custom soundfonts and samples that were already in the source .o2r file
	SkippedCentralDirSize         uint32 // Size in bytes of central directory entries for existing custom soundfonts and samples that were in the source .o2r file
	SkippedCount                  uint32 // Amount of existing custom soundfonts and samples that were in the source .o2r file
}

// Generate an asset map using mm.o2r and Audio.xml for .zbank -> Soundfont conversion. Due to this
// already having to read through the central directory of mm.o2r, this also picks out the read
// points for use when patching.
func NewAssetMap(o2rpath string, xmlpath string) (*Assets, *ReadPoints, error) {
	readPoints := ReadPoints{}
	fSrc, err := os.Open(o2rpath)
	if err != nil {
		return nil, nil, err
	}
	defer fSrc.Close()
	backwardsMap := map[string]uint32{}

	fSrc.Seek(-0xa, 2) // pointer to start of central directory
	var centralDirectoryLength uint32
	var centralDirectoryStart uint32
	binary.Read(fSrc, binary.LittleEndian, &centralDirectoryLength)
	binary.Read(fSrc, binary.LittleEndian, &centralDirectoryStart)
	readPoints.CentralDirectoryEnd = centralDirectoryStart + centralDirectoryLength
	fSrc.Seek(int64(centralDirectoryStart), 0)
	var headerSignature uint32
	var compressedSize uint32
	var filenameLength uint16
	var localHeaderOffset uint32
	if _, err := ReadUpToFilenameMatch(fSrc, &headerSignature, &compressedSize, &filenameLength, &localHeaderOffset, "audio/fonts/Soundfont_9"); err != nil {
		return nil, nil, err
	}
	sfCDirOffset, _ := fSrc.Seek(0x0, 1)
	readPoints.SoundfontsCentralDirectoryEnd = uint32(sfCDirOffset)
	readPoints.SoundfontsEnd = localHeaderOffset + 30 + uint32(filenameLength) + compressedSize
	if n, err := ReadUpToFilenameMatch(fSrc, &headerSignature, &compressedSize, &filenameLength, &localHeaderOffset, "audio/samples/Accordion_META"); err != nil {
		return nil, nil, err
	} else {
		readPoints.SkippedCount = n - 1
	}
	smpCDirOffset, _ := fSrc.Seek(0x0, 1)
	readPoints.SamplesCentralDirectoryBegin = uint32(smpCDirOffset) - 0x4A // the length of the Accordion central directory entry
	readPoints.SamplesBegin = localHeaderOffset
	readPoints.SkippedCentralDirSize = readPoints.SamplesCentralDirectoryBegin - readPoints.SoundfontsCentralDirectoryEnd
	readPoints.SkippedLocalSize = readPoints.SamplesBegin - readPoints.SoundfontsEnd
	SetSample(backwardsMap, "audio/samples/Accordion_META")
	for {
		filename, err := ReadCentralHeader(fSrc, &headerSignature, &compressedSize, &filenameLength, &localHeaderOffset)
		if err != nil {
			return nil, nil, fmt.Errorf("Reached end of archive too early\n")
		}
		SetSample(backwardsMap, filename)
		if filename == "audio/samples/sample_2_005412F0_META" {
			break
		}
	}
	readPoints.SamplesEnd = localHeaderOffset + 30 + uint32(filenameLength) + compressedSize
	sampleEntries, err := ParseXML(xmlpath)
	if err != nil {
		return nil, nil, err
	}
	finalMap := Assets{}
	for _, sample := range sampleEntries {
		currentOffset, _ := sample.GetOffset()
		if _, ok := backwardsMap[sample.Name]; ok {
			finalMap[currentOffset] = sample.Name + "_META"
		} else {
			finalMap[currentOffset] = fmt.Sprintf("sample_%s_%08X_META", sample.Bank, currentOffset)
		}
	}
	return &finalMap, &readPoints, nil
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
