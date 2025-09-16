package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev" // 由 build 时的 ldflags 设置

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  "显示 GoWeb CLI 的版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GoWeb CLI %s\n", version)
		fmt.Println("基于 GreenFish 的 Go Web 项目脚手架工具")
		fmt.Println("https://github.com/Programmer-Bugs-Bunny/GreenFishCli")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
