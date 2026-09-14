# Ship Packer
A custom sequence and instrument bank packaging tool for 2ship2harkinian and Ship of Harkinian

- [Features](#features)
- [What it does](#what-it-does)
- [Considerations](#considerations)
- [Setup](#setup)
    - [Z64 Song Packer](#z64-song-packer)
    - [Crashing songs](#songs-that-will-crash)
    - [Silent songs](#songs-that-dont-play)
- [Building](#building-from-source)
- [FAQ](#faq)

# Features

- [X] Support `.mmrs` sequences.
- [X] Support custom banks.
- [X] Support custom samples.
- [X] Support `.ootrs` sequences.
- \[Experimental\] Support for Ship of Harkinian (very crash prone with `.mmrs` sequences).

# What it does

This tool allows sequences that have custom instrument banks (`.mmrs` or `.ootrs` files that contain `.zbank` and `.bankmeta` files) to be played via 2ship2harkinian and Ship of Harkinian. It packages up all of the sequences you give it into a `.o2r` mod file. This tool accepts:

- `.ootrs` with and without custom banks
- `.mmrs` with and without custom banks
- `.*seq`, provided the name follows the "\[song\_name\]\_\[bank\_id\]\_\[category-list\].*seq" format
- `.seq`+`.meta` pairs, as you would give to retro

# Considerations

### 2Ship2Harkinian

While this can create .o2r files with essentially unlimited banks, the actual usable custom instrument **bank count is 214**. This limit will change in a future version of 2ship. As there is a bank limit, there is also a **sequence limit of 1919**. Ship Packer will give you warnings when you reach those limits.

Using multiple .o2r files with custom instrument banks will cause the sound fonts to overwrite each other, making sequences play with the wrong instruments.

To convert most `.ootrs` files for use in 2ship2harkinian, you must provide your copy `oot.o2r` in the same folder as your `shippacker` executable. This has only been tested so far using an `oot.o2r` that was generated from the **N64 NTSC 1.0 version** of the game. If you generated yours with a different version of the game, this might not work as intended. Further testing is required.

### Ship Of Harkinian

While this can generate .o2r's for Ship of Harkinian, it is very crash prone. More work needs to be done. Use this packer for that purpose with caution.

To convert most `.mmrs` files for use in Ship of Harkinian, you must provide your copy of `mm.o2r` in the same folder as your `shippacker` executable.

Category information is not currently preserved when packing for Ship of Harkinian.

# Setup

1) Download the [latest release](https://github.com/frogssoldseparately/shippacker/releases/latest) for your platform and unzip it to wherever you please.
2) Place whatever custom sequences you would like to pack inside the `music` folder.
3) If you would like to pack `.ootrs` sequences for 2ship2harkinian, place a copy of your `oot.o2r` in the same directory as your `shippacker` executable. If you would like to pack `.mmrs` sequences for Ship of Harkinian, place a copy of your `mm.o2r` in the same directory as your `shippacker` executable.
4) Run `shippacker.exe` and follow the terminal for further instructions. If all went well, your `mods` folder will house a file named `{some long number}.o2r`. You can rename it if you like. Move this file into 2ship2harkinian's (or Ship of Harkinian's) `mods` folder, boot it up, and your custom sequences will be readily available.

Any directory within the music folder whose name starts with an `_` (e.g., `_obscureGame/`) will not be explored by the packer. Consider these private folders. They could be useful if you already have a well organized tree of songs and wish to exclude certain folders (as well as its nested folders) from packing without having to delete anything.

## Z64 Song Packer

If you plan to use this in conjuction with the Z64 Song Packer, you have to export the songs as a `.zip` and extract the contents to use with this packer. Some songs might not work for you.

### Songs that *will* crash:

Darunia's Joy:

- Luigi's Mansion - Game Boy Horror.ootrs
- Mario & Luigi Superstar Saga - Come On! \[2\].ootrs
- Mario & Luigi Superstar Saga - Popple the Shadow Thief.ootrs
- Mario & Luigi Superstar Saga - The Last Cackletta.ootrs
- Pizza Tower - Boss Defeated.ootrs
- Pizza Tower - Tower Secret Treasure Found.ootrs
- Pokemon Black & White - Driftveil City.ootrs
- Sonic CD - Stardust Speedway Present (US).ootrs
- Sonic Spinball - Lava Powerhouse.ootrs
- The Legend of Zelda: Spirit Tracks - The Unenterable Body.ootrs
- Undertale Yellow - Mo Money.ootrs
- Yoshi's Island - Underground.ootrs

Japas' Jams:

- Touhou 15 - Pure Furies ~ Wherabouts of the Heart.mmrs

### Songs that don't play:

Darunia's Joy:

- Chrono Trigger - Manoria Cathedral.ootrs
- Final Fantasy VII - Electric de Chocobo.ootrs
- Guilty Gear X - Blue Water Blue Sky -May's Theme-.ootrs
- Kirby and the Forgotten Land - Running Through the New World.ootrs
- Mario Kart Double Dash!! - Rainbow Road.ootrs
- Mario & Luigi: Partners in Time - Battle.ootrs
- Metroid Fusion - Vs. Serris.ootrs
- The Legend of Zelda: Spirit Tracks - Fighting Cole and Malladus.ootrs
- Touhou 18 - Where is that Bustling Marketplace Now ~ Immemorial Marketeers.ootrs

Japas' Jams:

- Banjo-Tooie - Witchyworld \[2\].mmrs
- Chrono Trigger - Manoria Cathedral.mmrs
- Deltarune - Raise Up Your Bat.mmrs
- Duke Nukem 3D - The City Streets.mmrs
- Final Fantasy IX - Battle 2.mmrs
- Final Fantasy VI - Dancing Mad.mmrs
- Handel - Hallelujah Chorus.mmrs
- Holst - Mars, The Bringer of War - Part 1 (Looped).mmrs
- Memes - 4'33" (Silence).mmrs
- Omori - World's End Valentine.mmrs
- Old School Rune Scape - Test of Resourcefulness.mmrs
- Paper Mario Color Splash - World Map Medley.mmrs
- Pokemon Black 2 & White 2 - Battle! Colress.mmrs
- Pokemon Sun & Moon - Battle! Island Kahuna.mmrs
- Silent Hill 2 - Promise.mmrs
- Super Mario Galaxy - Peach's Castle is Stolen.mmrs
- Super Mario Odyssey - Steam Gardens.mmrs
- The Legend of Zelda: A Link to the Past - Magic Mirror.mmrs
- Touhou 8 - Reach for the Moon, Immortal Smoke.mmrs
- Touhou 13.5 - The Lost Emotion.mmrs
- Touhou 14 - Kobito of the Shining Needle ~ Little Princess.mmrs
- Touhou 15 - Eternal Spring Dream.mmrs
- Touhou 17 - Entrusting This World to Idols ~ Idolatrize World.mmrs
- Touhou 20 - Watatsuki's Spell Card ~ Divine Sea Battle.mmrs
- Turok 2 Seeds of Evil - Port of Adia.mmrs
- ULTRAKILL - ORDER (No Intro).mmrs
- ULTRAKILL - ORDER.mmrs

### Songs that' don't pack:

These won't make it into an `.o2r` using this tool, so they don't really need to be listed. There's 35 or so that won't pack because of invalid sample pointers (reused custom banks where the custom samples that aren't included are still referenced), or sometimes bad formatting, or something's wrong with my tool. These will be fixed eventually.

# Building from source

If you're on an unsupported platform, or would just like to build the executable yourself, please follow the steps below.

1) Install [Go](https://go.dev/doc/install).
2) Download the [shippacker](https://github.com/frogssoldseparately/shippacker) source code and unzip it.
3) Open the unpacked directory in a terminal and run the following commands:

On Windows:
```sh
go mod download
make native || go build -o ./bin/shippacker.exe ./cmd/shippacker/main.go
```

Anything else:
```sh
go mod download
make native
```

4) Move the `shippacker` file in the `bin` directory to anywhere of your choosing, or follow the rest of these steps from within the `bin` directory.
5) Create a `music` and `mods` folder in the folder that contains your `shippacker` file.
6) Follow [setup](#setup) from step 2 onwards.

# FAQ

## Q: Why do I need an mm/oot.o2r to pack certain songs?

A: While Ocarina of Time and Majora's Mask share a lot of assets, they don't share *all* of them. There are audio samples and soundfonts that are in one game but not the other. Your `mm.o2r` or `oot.o2r` is needed to get these missing assets so songs will play correctly on Ship of Harkinian and 2ship2harkinian respectively.

## Q: I have some songs that when played will play an entirely different song. Did something go wrong?

A: 2ship2harkinian has a limited amount of slots for sequences currently. Any songs over that limit will wrap around and play a song with a lower id number. If you completely packed an `.o2r` with 1919 songs, you cannot be using any other `.o2r`s that also add music. **Please note that it is not enough to exclude songs until you're back under this limit**. The `.o2r` itself needs to be at or under this limit. Shippacker gives warnings for this.

## Q: Some songs are playing with the wrong instruments or crashing. Can something be done?

A: If you are playing on Battler-Bravo or earlier, there is a limit of 255 total instrument banks. You might be using an `.o2r` that has more instrument banks than 2ship can handle, *or* you are using two different `.o2r`s that use custom banks at the same time. For Battler-Bravo and earlier, only one `.o2r` can include custom instrument banks. Otherwise, they overwrite each other.

## Q: I ensured I'm under both limits, I'm not using multiple music `.o2r`s that could be causing collisions, and yet there are still songs that don't play or crash. What else could be wrong?

A: While I've gotten almost all sequences of any kind to play without issue, there are still some that either [don't play](#songs-that-dont-play), will only play when you're in-game (not on the title screen), or [will just outright crash](#songs-that-will-crash). Some sequences could just be outright bad, or there's something I've failed to see thus far that makes them not play nice with the ports. I'm still figuring it out. Just exclude them for the time being. I apologize for not having a fix for them at the moment. If you run into one that isn't on the lists above, please let me know.