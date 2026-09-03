# Version tag. Defaults to something that will cause errors when omitted.
V=/
export V

build: windows linux-amd linux-arm mac-amd mac-arm wasm

windows: # Compile for Windows
	GOOS=windows GOARCH=amd64 go build -o ./bin/shippacker.exe ./cmd/shippacker/main.go

linux-amd: # Compile for Linux amd64
	GOOS=linux GOARCH=amd64 go build -o ./bin/shippacker_linux_amd ./cmd/shippacker/main.go

linux-arm: # Compile for Linux arm64
	GOOS=linux GOARCH=arm64 go build -o ./bin/shippacker_linux_arm ./cmd/shippacker/main.go

mac-amd: # Compile for Mac amd64
	GOOS=darwin GOARCH=amd64 go build -o ./bin/shippacker_mac_amd ./cmd/shippacker/main.go

mac-arm: # Compile for Mac arm64
	GOOS=darwin GOARCH=arm64 go build -o ./bin/shippacker_mac_arm ./cmd/shippacker/main.go

wasm: # Compile for WebAssembly
	GOOS=js GOARCH=wasm tinygo build -target=wasm -no-debug -o ./bin/shippacker.wasm ./cmd/shippacker/main_wasm.go

templates: # Generate pkg/globals files needed for compilation
	node utils/xmls/generateGlobalsFiles.js ./resources ./pkg/globals

dist: --dist-windows --dist-linux-amd --dist-linux-arm --dist-mac-amd --dist-mac-arm --dist-wasm # Generate distributables

--dist-windows: ./bin/shippacker.exe
# Prepare for release generation
	-mkdir ./dist
	mkdir ./dist/${V}
	mkdir ./dist/${V}/music
	mkdir ./dist/${V}/mods
# Remove possible duplicate zipped release only after version has been verified
# to be good so as to avoid any file system mishaps
	-rm "./dist/shippacker_${V}_windows_amd64.zip"
# Get Windows release resources
	cp ./bin/shippacker.exe ./dist/${V}/shippacker.exe
	cp ./README.md ./dist/${V}/README.md
# Make Windows release.
	cd ./dist/${V} ; 7z a "../shippacker_${V}_windows_amd64.zip" *
# Remove temporary uncompressed release folders. This could be a nasty command,
# but no one else should really be using it.
	rm -rf ./dist/${V}

--dist-linux-amd: ./bin/shippacker_linux_amd
# Prepare for release generation
	-mkdir ./dist
	mkdir ./dist/${V}
	mkdir ./dist/${V}/music
	mkdir ./dist/${V}/mods
# Remove possible duplicate zipped release only after version has been verified
# to be good so as to avoid any file system mishaps
	-rm "./dist/shippacker_${V}_linux_amd64.zip"
# Get Linux release resources
	cp ./bin/shippacker_linux_amd ./dist/${V}/shippacker
	cp ./README.md ./dist/${V}/README.md
# Make Linux release.
	cd ./dist/${V} ; 7z a "../shippacker_${V}_linux_amd64.zip" *
# Remove temporary uncompressed release folders.
	rm -rf ./dist/${V}

--dist-linux-arm: ./bin/shippacker_linux_arm
# Prepare for release generation
	-mkdir ./dist
	mkdir ./dist/${V}
	mkdir ./dist/${V}/music
	mkdir ./dist/${V}/mods
# Remove possible duplicate zipped release only after version has been verified
# to be good so as to avoid any file system mishaps
	-rm "./dist/shippacker_${V}_linux_arm64.zip"
# Get Linux release resources
	cp ./bin/shippacker_linux_arm ./dist/${V}/shippacker
	cp ./README.md ./dist/${V}/README.md
# Make Linux release.
	cd ./dist/${V} ; 7z a "../shippacker_${V}_linux_arm64.zip" *
# Remove temporary uncompressed release folders.
	rm -rf ./dist/${V}

--dist-mac-amd: ./bin/shippacker_mac_amd
# Prepare for release generation
	-mkdir ./dist
	mkdir ./dist/${V}
	mkdir ./dist/${V}/music
	mkdir ./dist/${V}/mods
# Remove possible duplicate zipped release only after version has been verified
# to be good so as to avoid any file system mishaps
	-rm "./dist/shippacker_${V}_mac_amd64.zip"
# Get Mac release resources
	cp ./bin/shippacker_mac_amd ./dist/${V}/shippacker
	cp ./README.md ./dist/${V}/README.md
# Make Mac release.
	cd ./dist/${V} ; 7z a "../shippacker_${V}_mac_amd64.zip" *
# Remove temporary uncompressed release folders.
	rm -rf ./dist/${V}

--dist-mac-arm: ./bin/shippacker_mac_arm
# Prepare for release generation
	-mkdir ./dist
	mkdir ./dist/${V}
	mkdir ./dist/${V}/music
	mkdir ./dist/${V}/mods
# Remove possible duplicate zipped release only after version has been verified
# to be good so as to avoid any file system mishaps
	-rm "./dist/shippacker_${V}_mac_arm64.zip"
# Get Mac release resources
	cp ./bin/shippacker_mac_arm ./dist/${V}/shippacker
	cp ./README.md ./dist/${V}/README.md
# Make Mac release.
	cd ./dist/${V} ; 7z a "../shippacker_${V}_mac_arm64.zip" *
# Remove temporary uncompressed release folders.
	rm -rf ./dist/${V}

--dist-wasm: ./bin/shippacker.wasm
	7z a "./dist/shippacker_${V}_js_wasm.zip" ./bin/shippacker.wasm