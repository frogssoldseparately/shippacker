package zbank

import (
	"github.com/frogssoldseparately/shippacker/pkg/sreader"
	"github.com/frogssoldseparately/shippacker/pkg/swriter"
)

func Read[T sreader.Number](r *sreader.SimpleReader) T {
	return sreader.Read[T](r)
}

func Write[T swriter.Number](w *swriter.SimpleWriter, data T) {
	swriter.Write(w, data)
}
