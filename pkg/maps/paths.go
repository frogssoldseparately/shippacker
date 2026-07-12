package maps

type Paths struct {
	MSrc string // Directory containing unmodified custom sequences
	OSrc string // Path to mm.o2r
	OOut string // Directory to write .o2r mod file in
	XSrc string // Path to Audio.xml file that comes with 2ship2harkinian
}

func BundlePaths(msrc string, osrc string, oout string, xsrc string) Paths {
	return Paths{msrc, osrc, oout, xsrc}
}
