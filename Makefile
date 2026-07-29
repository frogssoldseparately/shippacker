build:
	go build -o ./bin/shippacker.exe ./cmd/shippacker/main.go
xmls:
	node utils/packageXmls.js ./resources ./pkg/globals