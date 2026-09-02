# Ship Packer
A custom sequence and instrument bank packaging tool for 2ship2harkinian.

## Planned

- \[ In Progress \] Support .ootrs sequences.
- [X] Support custom samples.
- [ ] Support Ship of Harkinian.

## What it does

This tool allows sequences that have custom instrument banks (`.mmrs` files that contain `.zbank` and `.bankmeta` files) to be played via 2ship2harkinian. It packages up all of the sequences you give it into a `.o2r` mod file. This tool accepts:

- `.mmrs` with and without custom banks
- `.*seq`, provided the name follows the "\[song\_name\]\_\[bank\_id\]\_\[category-list\].*seq" format
- `.seq`+`.meta` pairs, as you would give to retro

## What it doesn't

This does not presently work for Ship of Harkinian, but support is planned.

## Considerations

While this can create .o2r files with essentially unlimited banks, the actual usable custom instrument bank count is 214. This limit will change in a future version of 2ship. As there is a bank limit, there is also a sequence limit of 1923. Ship Packer will give you warnings when you reach those limits.

Using multiple .o2r files with custom instrument banks will cause the sound fonts to overwrite each other, making sequences play with the wrong instruments.

## Setup
### Windows
1) Download the [latest release](https://github.com/frogssoldseparately/shippacker/releases/latest) and unzip it to wherever you please.
2) Place whatever custom sequences you would like to bundle into an `.o2r` file inside the `music` folder.
3) Run `shippacker.exe` and follow the terminal for further instructions. If all went well, your `mods` folder will house a file named `{some long number}.o2r`. You can rename it if you like. Move this file into 2ship2harkinian's `mods` folder, boot up 2ship2harkinian, and your custom sequences will be readily available.

### Building from source

If you're on a platform other than Windows, or would just like to build the executable yourself, please follow the steps below.

1) Install [Go](https://go.dev/doc/install).
2) Download the [shippacker](https://github.com/frogssoldseparately/shippacker/releases/latest) source code and unzip it.
3) Open the unpacked directory in a terminal and run the following commands:

On Windows:
```sh
go mod download
make windows || go build -o ./bin/shippacker.exe ./cmd/shippacker/main.go
```

On Linux amd64:
```sh
go mod download
make linux-amd
```

On Linux arm64:
```sh
go mod download
make linux-arm
```

For web deployment:
```sh
go mod download
make wasm # Requires tinygo
```

To generate distributables:
```sh
make build # Requires tinygo
make dist V=v[a version number] # Requires 7zip
```

4) Move the `shippacker.exe` (or `shippacker` on Linux) file in the `bin` directory to anywhere of your choosing, or follow the rest of these steps from within the `bin` directory.
5) Create a `music` and `mods` folder in the folder that contains your `shippacker` file.
6) Follow [setup](#windows) from step 2 onwards.