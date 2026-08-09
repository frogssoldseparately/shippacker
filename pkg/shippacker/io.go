package shippacker

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

func WriteRaw(w *swriter.SimpleWriter, buf *[]byte) {
	swriter.WriteRaw(w, buf)
}

func Read[T sreader.Number](r *sreader.SimpleReader) T {
	return sreader.Read[T](r)
}
