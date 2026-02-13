package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// seedFiles 种子文件与模块名的映射
var seedFiles = map[string]string{
	"permissions":   "001_seed_permissions.sql",
	"users":         "002_seed_users.sql",
	"stations":      "003_seed_stations.sql",
	"requests":      "004_seed_requests.sql",
	"notifications": "005_seed_notifications.sql",
	"news":          "006_seed_news.sql",
	"menus":         "007_seed_menus.sql",
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "执行数据库种子数据导入",
	Long: `执行数据库种子数据导入

示例:
  scare seed                    # 导入所有种子数据
  scare seed -t permissions     # 仅导入权限数据
  scare seed -t users,stations  # 导入多个模块
  scare seed --fresh            # 先 migrate --fresh 再导入`,
	Run: runSeed,
}

func init() {
	seedCmd.Flags().StringP("type", "t", "", "种子数据模块（逗号分隔）: permissions/users/stations/requests/notifications/news/menus")
	seedCmd.Flags().Bool("fresh", false, "先执行 migrate --fresh 再导入种子数据")
}

func runSeed(cmd *cobra.Command, args []string) {
	fresh, _ := cmd.Flags().GetBool("fresh")
	seedType, _ := cmd.Flags().GetString("type")

	// --fresh: 先执行 migrate --fresh
	if fresh {
		fmt.Println("🔄 --fresh: 先执行数据库迁移...")
		migrateCmd.Flags().Set("fresh", "true")
		migrateCmd.Flags().Set("auto-migrate", "true")
		runMigrate(migrateCmd, []string{})
	}

	_, db := initConfigAndDB(true)
	gdb := db.DB

	seedsDir := "database/seeds"

	// 确定要执行的种子文件列表
	var filesToExec []string
	if seedType != "" {
		// 按指定模块执行
		modules := strings.Split(seedType, ",")
		for _, m := range modules {
			m = strings.TrimSpace(m)
			filename, ok := seedFiles[m]
			if !ok {
				log.Fatalf("未知的种子模块: %s (可选: permissions/users/stations/requests/notifications/news/menus)", m)
			}
			filesToExec = append(filesToExec, filepath.Join(seedsDir, filename))
		}
	} else {
		// 执行所有 NNN_*.sql 文件（排除 seed_all.sql）
		files, err := filepath.Glob(filepath.Join(seedsDir, "[0-9][0-9][0-9]_*.sql"))
		if err != nil {
			log.Fatalf("读取种子目录失败: %v", err)
		}
		sort.Strings(files)
		filesToExec = files
	}

	if len(filesToExec) == 0 {
		fmt.Println("⚠️  没有找到种子文件")
		return
	}

	fmt.Printf("\n🌱 开始导入种子数据 (%d 个文件)...\n", len(filesToExec))

	success := 0
	for _, f := range filesToExec {
		filename := filepath.Base(f)
		content, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("读取种子文件 %s 失败: %v", filename, err)
		}

		sql := strings.TrimSpace(string(content))
		if sql == "" {
			continue
		}

		if err := gdb.Exec(sql).Error; err != nil {
			log.Fatalf("执行种子文件 %s 失败: %v", filename, err)
		}
		success++
		fmt.Printf("   ✅ %s\n", filename)
	}

	fmt.Printf("\n🎉 种子数据导入完成 (成功: %d/%d)\n", success, len(filesToExec))
}
