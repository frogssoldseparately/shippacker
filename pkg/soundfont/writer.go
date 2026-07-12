package soundfont

import "github.com/frogssoldseparately/shippacker/pkg/swriter"

func Write[T swriter.Number](w *swriter.SimpleWriter, data T) {
	swriter.Write(w, data)
}

func WriteString(w *swriter.SimpleWriter, s string, writeLength bool) {
	swriter.WriteString(w, s, writeLength)
}
