package service

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openServiceTestDB(t *testing.T, name string) *gorm.DB {
	t.Helper()

	tmpFile := t.TempDir() + "/" + name
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

	return db
}

func createUsersTable(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
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
			station_id INTEGER,
			status TEXT DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}
}

func createUserIdentitiesTable(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_identities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			identity_type TEXT,
			is_primary INTEGER,
			station_id INTEGER,
			status TEXT DEFAULT 'active',
			granted_at DATETIME,
			granted_by INTEGER,
			revoked_at DATETIME,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create user_identities table: %v", err)
	}
}

func createCustomerProfilesTable(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS customer_profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			id_card TEXT,
			address TEXT,
			latitude REAL,
			longitude REAL,
			customer_type TEXT,
			emergency_contact TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			gender TEXT,
			birth_date DATE,
			health_status TEXT,
			disability_level TEXT,
			medical_history TEXT,
			special_needs TEXT
		);
	`).Error; err != nil {
		t.Fatalf("failed to create customer_profiles table: %v", err)
	}
}

func createServiceRequestsTable(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS service_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			request_no TEXT,
			user_id INTEGER,
			service_type TEXT,
			status TEXT,
			description TEXT,
			submit_location_lat REAL,
			submit_location_lng REAL,
			service_location_lat REAL,
			service_location_lng REAL,
			contact_name TEXT,
			contact_phone TEXT,
			address TEXT,
			appointment_time DATETIME,
			urgency TEXT,
			source_station_id INTEGER,
			station_id INTEGER,
			dispatch_basis TEXT,
			needs_manual_verify INTEGER DEFAULT 0,
			reject_reason TEXT,
			images TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			rating INTEGER,
			feedback TEXT
		);
	`).Error; err != nil {
		t.Fatalf("failed to create service_requests table: %v", err)
	}
}

func createTaskAssignmentsTable(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS task_assignments (
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
}

func createServiceZonesTable(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS service_zones (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			station_id INTEGER,
			name TEXT,
			points TEXT,
			priority INTEGER,
			status TEXT DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
	`).Error; err != nil {
		t.Fatalf("failed to create service_zones table: %v", err)
	}
}
