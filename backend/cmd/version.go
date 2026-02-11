package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "1.0.0"
	BuildTime = "2026-02-07"
	GitCommit = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示 sCare 服务的版本信息`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sCare 社区养老服务平台\n")
		fmt.Printf("版本:    %s\n", Version)
		fmt.Printf("构建时间: %s\n", BuildTime)
		fmt.Printf("Git提交: %s\n", GitCommit)
	},
}
