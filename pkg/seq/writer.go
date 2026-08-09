package seq

import "github.com/frogssoldseparately/simpleseek/swriter"

func WriteRaw(w *swriter.SimpleWriter, buf *[]byte) {
	swriter.WriteRaw(w, buf)
}

func Write[T swriter.Number](w *swriter.SimpleWriter, data T) {
	swriter.Write(w, data)
}
