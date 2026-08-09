package maps

type Paths struct {
	MSrc string // Directory containing unmodified custom sequences
	OOut string // Directory to write .o2r mod file in
}

func BundlePaths(msrc string, oout string) Paths {
	return Paths{msrc, oout}
}
