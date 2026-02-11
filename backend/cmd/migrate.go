package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "执行数据库迁移",
	Long: `执行数据库迁移脚本

示例:
  scare migrate              # 执行所有迁移
  scare migrate --step 1     # 执行指定步骤`,
	Run: runMigrate,
}

func init() {
	migrateCmd.Flags().IntP("step", "s", 0, "执行指定步骤（0表示全部）")
}

func runMigrate(cmd *cobra.Command, args []string) {
	fmt.Println("数据库迁移功能开发中...")
	fmt.Println("请使用 SQL 脚本手动执行迁移: database/migrations/*.sql")
}
