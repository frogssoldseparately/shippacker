package swriter

type Int interface {
	int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

type Float interface {
	float32
}

type Number interface {
	Int | Float
}
