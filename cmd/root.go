// cmd/root.go

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd 作为包级变量定义
var rootCmd = &cobra.Command{
	Use:   "pngparser",
	Short: "A CLI tool to parse PNG files and display their chunks",
	Long: `PNGParser is a command-line application written in Go
that parses PNG files and displays information about each chunk,
including CRC verification.`,
	// 如果需要，可以在此添加 Run 函数
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 在此可以添加持久化标志或配置设置
}
