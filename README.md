# Ship Packer
A custom sequence and instrument bank packaging tool for 2ship2harkinian.

## What it does

This tool allows sequences that have custom instrument banks (`.mmrs` files that contain `.zbank` and `.bankmeta` files) to be played via 2ship2harkinian. It packages up all of the sequences you give it into a `.o2r` mod file that is placed directly in your 2ship `mods` folder (by default). This tool accepts:

- `.mmrs` with and without custom banks
- `.*seq`, provided the name follows the "\[song\_name\]\_\[bank\_id\]\_\[category-list\].*seq" format
- `.seq`+`.meta` pairs, as you would give to retro

## What it doesn't

This does not presently work for Ship of Harkinian, but support is planned.

This will ignore any sequences that have custom instruments (`.zsound` files), but support is planned.

## Setup
### Windows
1) Download the [latest release](https://github.com/frogssoldseparately/shippacker/releases/latest)
2) Place `shippacker.exe` in the directory that has your `2ship.exe`.

If you already have a folder that you like housing your custom sequences in, skip the rest of these steps and see [drag and drop](#drag-and-drop)

3) Create a folder named `music` in your `2ship.exe` directory. This is where you will put your custom sequences.
4) Place whatever custom sequences you would like to bundle into an `.o2r` file inside the `music` folder.
5) Double click `shippacker.exe` and note the terminal for any errors. If all went well, your `mods` folder will house a file named `{some long number}.o2r`. You can rename it if you like. Boot up 2ship2harkinian, and your custom sequences will be readily available.

### Drag and drop

Instead of creating a folder in your 2ship directory for music, you can reference your existing one through two methods. Via a terminal:

`.\shippacker.exe C:\absolute\path\to\music\folder`

Or you can open your 2ship folder and the folder that houses your music folder in the file explorer, and then drag the music folder onto `shippacker.exe`. Please note that you have to drag the folder, not the files within the folder.

After completing either of the above, your `.o2r` mod file containing custom sequences and banks will be located in your 2ship's `mods` folder.

For further options on usage, type the following into a terminal

`.\shippacker.exe -h`

### Building from source

If you're on a platform other than Windows, or would just like to build the executable yourself, please follow the steps below.

1) Install [Go](https://go.dev/doc/install).
2) Download the [shippacker](https://github.com/frogssoldseparately/shippacker/releases/latest) source code and unzip it.
3) Open the source directory in a terminal and run one of the following commands:

`make build`

or if you're on windows,

`go build -o ./bin/shippacker.exe ./cmd/shippacker/main.go`

4) Follow [setup](#windows) from step 2 onwards.