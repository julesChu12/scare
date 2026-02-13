package cmd

import (
	"fmt"
	"log"

	"community-elderly-care-platform/internal/config"
	"community-elderly-care-platform/pkg/database"
)

// initConfigAndDB 加载配置并初始化数据库连接
// multiStatements 为 true 时启用多语句执行（migrate/seed 需要）
func initConfigAndDB(multiStatements bool) (*config.Config, *database.DB) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	fmt.Printf("✅ 配置加载成功 (数据库: %s:%d/%s)\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	db, err := database.InitMySQL(database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		Charset:         cfg.Database.Charset,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MultiStatements: multiStatements,
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	fmt.Println("✅ 数据库连接成功")

	return cfg, db
}
