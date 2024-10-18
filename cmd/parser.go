// cmd/parse.go

package cmd

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Chunk represents a single PNG data chunk
type Chunk struct {
	Length uint32
	Type   string
	Data   []byte
	CRC    uint32
}

// Package-level variables for flags
var (
	outputDir       string
	maxDisplayLines int
	fixCRC          bool
)

// readChunk reads and parses a single chunk from PNG file
func readChunk(r io.Reader) (*Chunk, error) {
	var length uint32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	typeBytes := make([]byte, 4)
	if _, err := io.ReadFull(r, typeBytes); err != nil {
		return nil, err
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	var crc uint32
	if err := binary.Read(r, binary.BigEndian, &crc); err != nil {
		return nil, err
	}

	return &Chunk{
		Length: length,
		Type:   string(typeBytes),
		Data:   data,
		CRC:    crc,
	}, nil
}

// checkCRC checks the CRC of a chunk
func checkCRC(chunk *Chunk) bool {
	crc := crc32.NewIEEE()
	crc.Write([]byte(chunk.Type))
	crc.Write(chunk.Data)
	calculated := crc.Sum32()
	return calculated == chunk.CRC
}

// formatHex splits byte slice into a slice of hex strings, each representing a line with bytesPerLine bytes.
func formatHex(data []byte, bytesPerLine int) []string {
	hexData := hex.EncodeToString(data)
	var lines []string
	for i := 0; i < len(hexData); i += 2 * bytesPerLine {
		end := i + 2*bytesPerLine
		if end > len(hexData) {
			end = len(hexData)
		}
		line := hexData[i:end]
		var formattedLine string
		for j := 0; j < len(line); j += 2 {
			formattedLine += line[j:j+2] + " "
		}
		lines = append(lines, formattedLine)
	}
	return lines
}

// computeCRC32 computes the CRC32 for the given chunk type and data
func computeCRC32(chunkType string, data []byte) uint32 {
	crc := crc32.NewIEEE()
	crc.Write([]byte(chunkType))
	crc.Write(data)
	return crc.Sum32()
}

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse [file]",
	Short: "Parse a PNG file and display its chunks",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]

		// 获取标志值
		verbose, _ := cmd.Flags().GetBool("verbose")
		// outputDir 和 maxDisplayLines 已在包级别声明，无需再次声明
		// fixCRC 也在包级别声明
		// 变量已通过 StringVarP 和 IntVarP 绑定
		// 直接使用 outputDir, maxDisplayLines, fixCRC

		// 确保输出目录存在
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			fmt.Printf("Output directory does not exist: %s\n", outputDir)
			return
		}

		// 以读写模式打开文件，以便在需要时进行修复
		file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		// Read the PNG signature
		signature := make([]byte, 8)
		if _, err := io.ReadFull(file, signature); err != nil {
			fmt.Printf("Error reading PNG signature: %v\n", err)
			return
		}

		if !bytes.Equal(signature, []byte("\x89PNG\r\n\x1a\n")) {
			fmt.Println("Not a valid PNG file")
			return
		}

		fmt.Println("Valid PNG file detected. Parsing chunks...")

		var offset int64 = 8 // PNG signature 已读 8 字节

		// Read chunks
		for {
			chunk, err := readChunk(file)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("Error reading chunk at offset %d: %v\n", offset, err)
				return
			}

			fmt.Printf("\nChunk Type: %s, Length: %d bytes, Offset: %d - %d\n", chunk.Type, chunk.Length, offset, offset+int64(12+chunk.Length))
			// 12 = 4 (Length) + 4 (Type) + 4 (CRC)

			// 检查 CRC
			if checkCRC(chunk) {
				fmt.Println("CRC: 正确")
			} else {
				fmt.Println("CRC: 错误")

				// 如果是 IHDR 块，提供修复建议
				if chunk.Type == "IHDR" {
					fmt.Println("检测到 IHDR 块的 CRC 错误。")
					correctCRC := computeCRC32(chunk.Type, chunk.Data)
					fmt.Printf("建议修复：将 CRC 修改为 0x%08X\n", correctCRC)
					fmt.Printf("您可以使用十六进制编辑器（如 hexedit、Hex Fiend 等）打开文件并将 CRC 值更新为 0x%08X。\n", correctCRC)

					// 如果启用了自动修复
					if fixCRC {
						fmt.Printf("自动修复 IHDR CRC: 0x%08X\n", correctCRC)

						// CRC 位于该 Chunk 的最后 4 字节
						crcPosition := offset + int64(8+chunk.Length)
						_, err := file.Seek(crcPosition, io.SeekStart)
						if err != nil {
							fmt.Printf("Error seeking to CRC position: %v\n", err)
							return
						}

						// 写入正确的 CRC32 值（大端字节序）
						crcBytes := make([]byte, 4)
						binary.BigEndian.PutUint32(crcBytes, correctCRC)
						_, err = file.Write(crcBytes)
						if err != nil {
							fmt.Printf("Error writing correct CRC: %v\n", err)
							return
						}
						fmt.Println("CRC 自动修复完成。")
					}
				}

				// 格式化 Hex 数据，每行 16 字节
				lines := formatHex(chunk.Data, 16)

				// 定义显示的最大行数
				if len(lines) > maxDisplayLines {
					// 生成文件名，避免冲突
					safeType := strings.TrimSpace(chunk.Type)
					filename := fmt.Sprintf("chunk_%s_0x%X.hex", safeType, offset)
					filename = filepath.Join(outputDir, filepath.Clean(filename))

					// 保存所有 Hex 数据到文件
					fullHexData := strings.Join(lines, "\n")
					err := os.WriteFile(filename, []byte(fullHexData), 0644)
					if err != nil {
						fmt.Printf("Error saving hex data to file: %v\n", err)
					} else {
						// 打印前 maxDisplayLines 行
						fmt.Printf("Chunk Data (Hex) [First %d lines]:\n", maxDisplayLines)
						for i := 0; i < maxDisplayLines && i < len(lines); i++ {
							fmt.Println(lines[i])
						}
						fmt.Printf("Full hex data saved to file: %s\n", filename)
					}
				} else {
					// 数据行数不超过限制，直接打印所有
					if verbose {
						fmt.Println("Chunk Data (Hex):")
						for _, line := range lines {
							fmt.Println(line)
						}
					} else {
						fmt.Println("Chunk Data (Hex):")
						for _, line := range lines {
							fmt.Println(line)
						}
					}
				}
			}

			// 更新偏移量
			offset += int64(12 + chunk.Length)
		}
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	// 添加 --verbose 标志
	parseCmd.Flags().BoolP("verbose", "v", false, "Show detailed hex data for chunks with CRC errors")

	// 添加 --output-dir 标志
	parseCmd.Flags().StringVarP(&outputDir, "output-dir", "o", ".", "Directory to save hex data files")

	// 添加 --max-lines 标志
	parseCmd.Flags().IntVarP(&maxDisplayLines, "max-lines", "m", 10, "Number of hex lines to display before saving to file")

	// 添加 --fix 标志
	parseCmd.Flags().BoolVarP(&fixCRC, "fix", "f", false, "Automatically fix CRC errors in IHDR chunk")
}
