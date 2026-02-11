// gorm_gen.go - GORM Gen tool for generating model and query files
package main

import (
	"flag"
	"fmt"
	"os"

	"community-elderly-care-platform/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./configs", "配置文件目录路径")
	flag.Parse()

	fmt.Println("Starting GORM model generation...")

	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// 2. 构建 DSN 并连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully connected to the database.")

	// 3. 初始化生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./internal/dao/query",
		ModelPkgPath:      "./internal/dao/model",
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldNullable:     false,
		FieldCoverable:    false,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})
	g.UseDB(db)

	// 4. 获取所有表
	tables, err := db.Migrator().GetTables()
	if err != nil {
		fmt.Printf("Failed to get tables: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d tables in database\n", len(tables))

	// 5. 打印表信息
	var filteredTables []string
	for _, table := range tables {
		filteredTables = append(filteredTables, table)
		fmt.Printf("Including table: %s\n", table)
	}

	fmt.Printf("Generating models and queries for %d filtered tables...\n", len(filteredTables))

	// 6. 为过滤后的表生成模型
	var allModels []interface{}
	for _, tableName := range filteredTables {
		model := g.GenerateModel(tableName)
		allModels = append(allModels, model)
	}

	// 7. 应用基础查询方法
	g.ApplyBasic(allModels...)

	// 8. 执行生成
	g.Execute()

	fmt.Println("GORM model and query generation finished successfully.")
}
