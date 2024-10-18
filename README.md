# 项目介绍

本项目提供给ctf 图片分析

```
picture_data_check -h
PNGParser is a command-line application written in Go
that parses PNG files and displays information about each chunk,
including CRC verification.

Usage:
  pngparser [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  parse       Parse a PNG file and display its chunks

Flags:
  -h, --help   help for pngparser

Use "pngparser [command] --help" for more information about a command.
```

```
picture_data_check  parse -h
Parse a PNG file and display its chunks

Usage:
  pngparser parse [file] [flags]

Flags:
  -f, --fix                 Automatically fix CRC errors in IHDR chunk
  -h, --help                help for parse
  -m, --max-lines int       Number of hex lines to display before saving to file (default 10)
  -o, --output-dir string   Directory to save hex data files (default ".")
  -v, --verbose             Show detailed hex data for chunks with CRC errors

```
使用样例：
```shell
picture parse png\1.PNG
```
```
❯ picture_data_check.exe  parse png\1.PNG
Valid PNG file detected. Parsing chunks...

Chunk Type: IHDR, Length: 13 bytes, Offset: 8 - 33
CRC: 错误
检测到 IHDR 块的 CRC 错误。
建议修复：将 CRC 修改为 0x01B070E5
您可以使用十六进制编辑器（如 hexedit、Hex Fiend 等）打开文件并将 CRC 值更新为 0x01B070E5。
Chunk Data (Hex):
d5 05 2b fc d5 05 2b fc 08 06 00 00 00

Chunk Type: IDAT, Length: 65445 bytes, Offset: 33 - 65490
CRC: 正确

Chunk Type: IDAT, Length: 15717 bytes, Offset: 65490 - 81219
CRC: 正确

Chunk Type: IEND, Length: 0 bytes, Offset: 81219 - 81231
CRC: 正确

```