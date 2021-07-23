[![build](https://github.com/ahmadkabdullah/winela/actions/workflows/build-test.yml/badge.svg?branch=main)](https://github.com/ahmadkabdullah/winela/actions/workflows/build-test.yml)
![license](https://img.shields.io/github/license/ahmadkabdullah/winela?label=License&style=flat&color=yellow)
![lines](https://img.shields.io/tokei/lines/github/ahmadkabdullah/winela?label=Lines)

Winela is a commandline launcher for executables through winehq.

## Usage
Winela operates on two files stored in **winela** dir which in turn is in config dir (usually **~/.config/winela/**):
- **wineladb**: storing list of exes to launch
- **winelarc**: containing configuration for specifying wine version and parameters.

You can *scan* to populate **wineladb** with exe files in a directory (on first run, this also creates the other file).

You can *list* the exe files acquired from the scan. This reads out a numerated version of **wineladb**.

You can *run* an item from the (numerated) list. Also can choose to fork the process or not.
