package soundfont

import (
	"github.com/frogssoldseparately/simpleseek/sreader"
	"github.com/frogssoldseparately/simpleseek/swriter"
)

func Write[T swriter.Number](w *swriter.SimpleWriter, data T) {
	swriter.Write(w, data)
}

func WriteString(w *swriter.SimpleWriter, s string, writeLength bool) {
	swriter.WriteString(w, s, writeLength)
}

func Read[T sreader.Number](r *sreader.SimpleReader) T {
	return sreader.Read[T](r)
}

func ReadString(r *sreader.SimpleReader, len uint32) string {
	buf := []byte{}
	for range len {
		buf = append(buf, Read[uint8](r))
	}
	return string(buf)
}
