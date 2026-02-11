package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scare",
	Short: "sCare 社区养老服务平台",
	Long: `sCare 社区养老服务平台后端服务

一个基于 Go + Gin + GORM 的社区养老服务信息分发平台。
支持 B端管理和 C端用户服务。`,
}

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(seedCmd)
	rootCmd.AddCommand(versionCmd)
}
