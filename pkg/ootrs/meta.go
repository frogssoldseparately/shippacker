package ootrs

import (
	"archive/zip"
	"io"
	"strings"
)

type OotrsMeta struct {
	Name          string
	Bank          string
	Type          string
	Categories    []string
	CustomSamples []string
}

func processOotrsMeta(f *zip.File) (*OotrsMeta, error) {
	r, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	rawLines := strings.Split(string(buf), "\n")
	lines := []string{}
	zsoundLines := []string{}
	for _, rawLine := range rawLines {
		line := strings.Trim(rawLine, " \r\t")
		if len(line) > 0 {
			if strings.Contains(line, "ZSOUND") {
				zsoundLines = append(zsoundLines, line)
			} else if len(lines) == 4 {
				lines[3] = line
			} else {
				lines = append(lines, line)
			}
		}
	}
	if len(lines) == 2 {
		lines = append(lines, "bgm")
	}
	if len(lines) == 3 {
		lines = append(lines, "")
	}
	return &OotrsMeta{
		lines[0],
		lines[1],
		lines[2],
		strings.FieldsFunc(lines[3], func(r rune) bool { return r == ',' || r == '-' }),
		zsoundLines,
	}, nil
}
