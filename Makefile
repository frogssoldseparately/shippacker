build:
	go build -o ./bin/shippacker.exe ./cmd/shippacker/main.go
xmls:
	node utils/xmls/packageAudioXmls.js ./resources ./pkg/globals