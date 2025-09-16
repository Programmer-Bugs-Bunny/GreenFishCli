package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goweb",
	Short: "Go Web 项目脚手架工具",
	Long: `goweb 是一个基于 GoWebTemplate 的项目脚手架生成工具。
可以快速创建标准化的 Go Web 项目，包含 Gin、GORM、Zap 日志等常用组件。`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "详细输出")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
