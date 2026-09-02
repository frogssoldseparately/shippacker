package o2r

const (
	LittleEndian uint32 = iota
	BigEndian
)

const (
	SequenceType    = 0x4F534551
	SequenceVersion = 0x0
	SequenceCustom  = 0x1

	SequenceMedium = 0x2
)

const (
	SoundfontType    = 0x4F534654
	SoundfontVersion = 0x0
	SoundfontCustom  = 0x1
)

const (
	SampleType    = 0x4F534D50
	SampleVersion = 0x0
	SampleCustom  = 0x0
)

const (
	Attributes = 0x81B60000
)

const (
	Store        = 0x0
	StoreFlags   = 0x800
	Deflate      = 0x8
	DeflateFlags = 0x802
)

type Resource interface {
	GetEndianness() uint32
	GetResourceType() uint32
	GetResourceVersion() uint32
	GetIsCustom() uint8
}
