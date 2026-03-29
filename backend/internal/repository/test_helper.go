package repository

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试用的 SQLite 数据库
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	tmpFile := t.TempDir() + "/test.db"
	dsn := tmpFile + "?_loc=auto&parseTime=true"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NowFunc: func() time.Time {
			return time.Now().Local().Truncate(time.Second)
		},
	})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// 创建所有需要的表
	createTables(t, db)

	return db
}

// createTables 创建测试所需的表结构
func createTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	// users 表
	if err := db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			phone TEXT UNIQUE,
			password_hash TEXT,
			name TEXT,
			email TEXT,
			avatar TEXT,
			gender TEXT,
			birth_date DATE,
			id_card TEXT,
			id_card_hmac TEXT,
			id_card_masked TEXT,
			status TEXT DEFAULT 'active',
			station_id INTEGER,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	// user_identities 表
	if err := db.Exec(`
		CREATE TABLE user_identities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			identity_type TEXT,
			status TEXT DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create user_identities table: %v", err)
	}

	// service_requests 表
	if err := db.Exec(`
		CREATE TABLE service_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			request_no TEXT UNIQUE,
			user_id INTEGER,
			service_type TEXT,
			status TEXT,
			description TEXT,
			submit_location_lat REAL,
			submit_location_lng REAL,
			contact_name TEXT,
			contact_phone TEXT,
			address TEXT,
			appointment_time DATETIME,
			urgency TEXT,
			station_id INTEGER,
			reject_reason TEXT,
			images TEXT,
			rating INTEGER,
			feedback TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create service_requests table: %v", err)
	}

	// task_assignments 表
	if err := db.Exec(`
		CREATE TABLE task_assignments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			request_id INTEGER,
			station_id INTEGER,
			staff_id INTEGER,
			status TEXT,
			claimed_at DATETIME,
			completed_at DATETIME,
			rating INTEGER,
			feedback TEXT,
			staff_notes TEXT,
			images TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create task_assignments table: %v", err)
	}

	// service_stations 表
	if err := db.Exec(`
		CREATE TABLE service_stations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			code TEXT,
			address TEXT,
			phone TEXT,
			latitude REAL,
			longitude REAL,
			service_area TEXT,
			capacity INTEGER DEFAULT 10,
			work_hours TEXT,
			status TEXT DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create service_stations table: %v", err)
	}

	// service_zones 表
	if err := db.Exec(`
		CREATE TABLE service_zones (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			station_id INTEGER,
			name TEXT,
			points TEXT,
			priority INTEGER DEFAULT 0,
			status TEXT DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create service_zones table: %v", err)
	}

	// notifications 表
	if err := db.Exec(`
		CREATE TABLE notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			title TEXT,
			body TEXT,
			type TEXT,
			related_id INTEGER,
			related_type TEXT,
			channel TEXT,
			send_status TEXT DEFAULT 'pending',
			sent_at DATETIME,
			is_read INTEGER DEFAULT 0,
			read_at DATETIME,
			retry_count INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create notifications table: %v", err)
	}

	// customer_profiles 表
	if err := db.Exec(`
		CREATE TABLE customer_profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			customer_type TEXT,
			gender TEXT,
			address TEXT,
			latitude REAL,
			longitude REAL,
			emergency_contact TEXT,
			health_status TEXT,
			disability_level TEXT,
			medical_history TEXT,
			special_needs TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create customer_profiles table: %v", err)
	}
}
