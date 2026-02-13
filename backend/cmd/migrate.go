package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"community-elderly-care-platform/internal/dao/model"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "执行数据库迁移",
	Long: `执行数据库迁移（GORM AutoMigrate + SQL 增量迁移）

示例:
  scare migrate                # 执行 AutoMigrate + SQL 迁移
  scare migrate --auto-migrate=false  # 仅执行 SQL 迁移
  scare migrate --fresh        # 删除所有表后重新迁移`,
	Run: runMigrate,
}

func init() {
	migrateCmd.Flags().Bool("auto-migrate", true, "执行 GORM AutoMigrate")
	migrateCmd.Flags().Bool("fresh", false, "删除所有表后重新迁移（危险操作）")
}

// allModels 返回所有需要 AutoMigrate 的模型
func allModels() []interface{} {
	return []interface{}{
		&model.User{},
		&model.UserIdentity{},
		&model.Role{},
		&model.Permission{},
		&model.RolePermission{},
		&model.Menu{},
		&model.ServiceStation{},
		&model.ServiceZone{},
		&model.CustomerProfile{},
		&model.ServiceRequest{},
		&model.TaskAssignment{},
		&model.TaskHistory{},
		&model.Notification{},
		&model.News{},
		&model.Banner{},
		&model.Report{},
	}
}

// SchemaMigration 迁移记录表
type SchemaMigration struct {
	ID        uint   `gorm:"primaryKey"`
	Filename  string `gorm:"size:255;uniqueIndex"`
	AppliedAt time.Time
}

func (SchemaMigration) TableName() string {
	return "schema_migrations"
}

func runMigrate(cmd *cobra.Command, args []string) {
	fresh, _ := cmd.Flags().GetBool("fresh")
	autoMigrate, _ := cmd.Flags().GetBool("auto-migrate")

	_, db := initConfigAndDB(true)
	gdb := db.DB

	// --fresh: 删除所有表后重新迁移
	if fresh {
		fmt.Println("⚠️  --fresh: 删除所有表...")
		tables, err := getTables(gdb)
		if err != nil {
			log.Fatalf("获取表列表失败: %v", err)
		}
		if len(tables) > 0 {
			// 禁用外键检查
			gdb.Exec("SET FOREIGN_KEY_CHECKS = 0")
			for _, t := range tables {
				gdb.Exec("DROP TABLE IF EXISTS `" + t + "`")
				fmt.Printf("   已删除: %s\n", t)
			}
			gdb.Exec("SET FOREIGN_KEY_CHECKS = 1")
		}
		fmt.Println("✅ 所有表已删除")
	}

	// Step 1: GORM AutoMigrate
	if autoMigrate {
		fmt.Println("\n📦 Step 1: GORM AutoMigrate...")
		models := allModels()
		if err := gdb.AutoMigrate(models...); err != nil {
			log.Fatalf("AutoMigrate 失败: %v", err)
		}
		fmt.Printf("✅ AutoMigrate 完成 (%d 个模型)\n", len(models))
	}

	// Step 2: 执行增量 SQL 迁移
	fmt.Println("\n📦 Step 2: 执行增量 SQL 迁移...")
	// 确保 schema_migrations 表存在
	if err := gdb.AutoMigrate(&SchemaMigration{}); err != nil {
		log.Fatalf("创建 schema_migrations 表失败: %v", err)
	}

	migrationsDir := "database/migrations"
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		log.Fatalf("读取迁移目录失败: %v", err)
	}
	sort.Strings(files)

	applied := 0
	skipped := 0
	for _, f := range files {
		filename := filepath.Base(f)

		// 检查是否已执行
		var count int64
		gdb.Model(&SchemaMigration{}).Where("filename = ?", filename).Count(&count)
		if count > 0 {
			skipped++
			continue
		}

		// 读取并执行 SQL
		content, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("读取迁移文件 %s 失败: %v", filename, err)
		}

		sql := strings.TrimSpace(string(content))
		if sql == "" {
			continue
		}

		if err := gdb.Exec(sql).Error; err != nil {
			log.Fatalf("执行迁移 %s 失败: %v", filename, err)
		}

		// 记录已执行
		gdb.Create(&SchemaMigration{
			Filename:  filename,
			AppliedAt: time.Now(),
		})
		applied++
		fmt.Printf("   ✅ %s\n", filename)
	}

	fmt.Printf("\n🎉 迁移完成 (执行: %d, 跳过: %d)\n", applied, skipped)
}

// getTables 获取数据库中所有表名
func getTables(db *gorm.DB) ([]string, error) {
	var tables []string
	if err := db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}
