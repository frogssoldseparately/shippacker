# Ship Packer
A custom sequence and instrument bank packaging tool for 2ship2harkinian.

## What it does

This tool allows sequences that have custom instrument banks (`.mmrs` files that contain `.zbank` and `.bankmeta` files) to be played via 2ship2harkinian. It packages up all of the sequences you give it into a `.o2r` mod file. This tool accepts:

- `.mmrs` with and without custom banks
- `.*seq`, provided the name follows the "\[song\_name\]\_\[bank\_id\]\_\[category-list\].*seq" format
- `.seq`+`.meta` pairs, as you would give to retro

## What it doesn't

This does not presently work for Ship of Harkinian, but support is planned.

This will ignore any sequences that have custom instruments (`.zsound` files), but support is planned.

## Setup
### Windows
1) Download the [latest release](https://github.com/frogssoldseparately/shippacker/releases/latest) and unzip it to wherever you please.
2) Place whatever custom sequences you would like to bundle into an `.o2r` file inside the `music` folder.
3) Run `shippacker.exe` and note the terminal for any errors. If all went well, your `mods` folder will house a file named `{some long number}.o2r`. You can rename it if you like. Move this file into 2ship2harkinian's `mods` folder, boot up 2ship2harkinian, and your custom sequences will be readily available.

### Building from source

If you're on a platform other than Windows, or would just like to build the executable yourself, please follow the steps below.

1) Install [Go](https://go.dev/doc/install).
2) Download the [shippacker](https://github.com/frogssoldseparately/shippacker/releases/latest) source code and unzip it.
3) Open the unpacked directory in a terminal and run the following commands:

On Windows:
```sh
go mod download
make build-w || go build -o ./bin/shippacker.exe ./cmd/shippacker/main.go
```

On Linux:
```sh
go mod download
make build-l # L for linux
```

4) Move the `shippacker.exe` (or `shippacker` on Linux) file in the `bin` directory to anywhere of your choosing, or follow the rest of these steps from within the `bin` directory.
5) Create a `music` and `mods` folder in the folder that contains your `shippacker` file.
6) Follow [setup](#windows) from step 2 onwards.