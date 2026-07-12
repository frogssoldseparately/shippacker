package o2r

import "github.com/frogssoldseparately/shippacker/pkg/swriter"

func Write[T swriter.Number](w *swriter.SimpleWriter, data T) {
	swriter.Write(w, data)
}

func WriteHeader(resource Resource, w *swriter.SimpleWriter) {
	Write(w, uint32(resource.GetEndianness()))
	Write(w, resource.GetResourceType())
	Write[uint32](w, 0x2)
	Write[uint32](w, 0xDEADBEEF)
	Write[uint32](w, 0xDEADBEEF)
	Write(w, resource.GetResourceVersion())
	Write(w, resource.GetIsCustom())
	Write[uint8](w, 0x0)
	Write[uint8](w, 0x0)
	Write[uint8](w, 0x0)
	Write[uint32](w, 0x0)
	Write[uint32](w, 0x0)
	for w.GetLength() < 0x40 {
		Write[uint32](w, 0x0)
	}
}
