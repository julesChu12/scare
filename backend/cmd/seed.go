package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "执行数据库种子数据导入",
	Long: `执行数据库种子数据导入

示例:
  scare seed              # 导入所有种子数据
  scare seed --type users # 导入指定类型`,
	Run: runSeed,
}

func init() {
	seedCmd.Flags().StringP("type", "t", "", "种子数据类型: users/stations/permissions/news")
}

func runSeed(cmd *cobra.Command, args []string) {
	fmt.Println("种子数据导入功能开发中...")
	fmt.Println("请使用 SQL 脚本手动执行导入: database/seeds/*.sql")
}
