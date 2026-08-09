# Version tag. Defaults to something that will cause errors when omitted.
V=/
export V

build-w: # Compile for Windows
	go build -o ./bin/shippacker.exe ./cmd/shippacker/main.go

build-l: # Compile for Linux?
	go build -o ./bin/shippacker ./cmd/shippacker/main.go

xmls: # Generate audio xml files needed for compilation
	node utils/xmls/packageAudioXmls.js ./resources ./pkg/globals

dist: ./bin/shippacker.exe # Generate distributables
# Prepare for release generation
	-mkdir ./dist
	mkdir ./dist/${V}
	mkdir ./dist/${V}/music
	mkdir ./dist/${V}/mods
# Remove possible duplicate zipped release only after version has been verified
# to be good so as to avoid any file system mishaps
	-rm "./dist/Ship Packer ${V}.zip"
# Get Windows release resources
	cp ./bin/shippacker.exe ./dist/${V}/shippacker.exe
	cp ./README.md ./dist/${V}/README.md
# Make Windows release. cd'ing in ensures the outermost folder isn't also zipped
	cd ./dist/${V} ; 7z a "../Ship Packer ${V}.zip" *
# Remove temporary uncompressed release folders. This could be a nasty command,
# but no one else should really be using it.
	rm -rf ./dist/${V}