# Ship Packer
A custom sequence and instrument bank packaging tool for 2ship2harkinian.

## Planned

- [ ] Support `.ootrs` sequences.
    - [X] Convert `.ootrs` to 2ship2harkinian's format.
    - [ ] Change sample injection infrastructure so bank 0 can be supported.
    - [ ] Convert `.ootrs` categories to `.mmrs` categories.
    - [ ] Better parse Ship of Harkinian soundfonts.
- [X] Support custom samples.
- [ ] Support for Ship of Harkinian o2r.

## What it does

This tool allows sequences that have custom instrument banks (`.mmrs` or `.ootrs` files that contain `.zbank` and `.bankmeta` files) to be played via 2ship2harkinian. It packages up all of the sequences you give it into a `.o2r` mod file. This tool accepts:

- `.ootrs` with and without custom banks
- `.mmrs` with and without custom banks
- `.*seq`, provided the name follows the "\[song\_name\]\_\[bank\_id\]\_\[category-list\].*seq" format
- `.seq`+`.meta` pairs, as you would give to retro

## What it doesn't

This does not presently work for Ship of Harkinian, but support is planned.

## Considerations

While this can create .o2r files with essentially unlimited banks, the actual usable custom instrument **bank count is 214**. This limit will change in a future version of 2ship. As there is a bank limit, there is also a **sequence limit of 1923**. Ship Packer will give you warnings when you reach those limits.

Using multiple .o2r files with custom instrument banks will cause the sound fonts to overwrite each other, making sequences play with the wrong instruments.

To convert most `.ootrs` files, you must provide your copy `oot.o2r` in the same folder as your `shippacker` executable. This has only been tested so far using an `oot.o2r` that was generated from the **N64 NTSC 1.0 version** of the game. If you generated yours with a different version of the game, this might not work as intended. Further testing is required.

Due to how the sample injections work currently, any `.ootrs` using instrument bank 0 will be ignored.

Category information is currently not preserved when including an `.ootrs` file. It defaults to either `bgm` or `fanfare`.

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

On Mac amd64:
```sh
go mod download
make mac-amd
```

On Mac arm64:
```sh
go mod download
make mac-arm
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

4) Move the `shippacker` file in the `bin` directory to anywhere of your choosing, or follow the rest of these steps from within the `bin` directory.
5) Create a `music` and `mods` folder in the folder that contains your `shippacker` file.
6) Follow [setup](#windows) from step 2 onwards.