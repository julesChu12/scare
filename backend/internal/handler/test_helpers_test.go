package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"community-elderly-care-platform/internal/dao/model"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func openHandlerTestDB(t *testing.T, name string) *gorm.DB {
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

func createHandlerTables(t *testing.T, db *gorm.DB) {
	t.Helper()

	statements := []string{
		`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			phone TEXT UNIQUE,
			password_hash TEXT NOT NULL,
			name TEXT,
			email TEXT,
			avatar TEXT,
			gender TEXT,
			birth_date DATE,
			id_card TEXT,
			id_card_hmac TEXT,
			id_card_masked TEXT,
			station_id INTEGER DEFAULT 0,
			status TEXT DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE user_identities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			identity_type TEXT NOT NULL,
			is_primary INTEGER NOT NULL DEFAULT 0,
			station_id INTEGER DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'active',
			granted_at DATETIME,
			granted_by INTEGER,
			revoked_at DATETIME,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE customer_profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
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
		`,
		`
		CREATE TABLE service_stations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
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
		`,
		`
		CREATE TABLE service_zones (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			station_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			points TEXT NOT NULL,
			priority INTEGER DEFAULT 0,
			status TEXT DEFAULT 'active',
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
		`
		CREATE TABLE service_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			request_no TEXT,
			user_id INTEGER NOT NULL,
			service_type TEXT NOT NULL,
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
			source_station_id INTEGER DEFAULT 0,
			station_id INTEGER DEFAULT 0,
			dispatch_basis TEXT,
			needs_manual_verify INTEGER DEFAULT 0,
			reject_reason TEXT,
			images TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME,
			rating INTEGER DEFAULT 0,
			feedback TEXT,
			station_name TEXT,
			source_station_name TEXT
		);
		`,
		`
		CREATE TABLE task_assignments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			request_id INTEGER NOT NULL,
			station_id INTEGER NOT NULL,
			staff_id INTEGER DEFAULT 0,
			status TEXT,
			claimed_at DATETIME,
			completed_at DATETIME,
			rating INTEGER DEFAULT 0,
			feedback TEXT,
			staff_notes TEXT,
			images TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		`,
	}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to create tables: %v", err)
		}
	}
}

func newJSONTestContext(t *testing.T, method, path string, payload any) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	var body *bytes.Buffer
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		body = bytes.NewBuffer(raw)
	} else {
		body = bytes.NewBuffer(nil)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, w
}

func setBEndClaims(c *gin.Context, userID, stationID int64, identities ...string) {
	c.Set("user_id", userID)
	c.Set("user_type", "b_end")
	c.Set("station_id", stationID)
	c.Set("user_identities", identities)
}

func setCEndClaims(c *gin.Context, userID int64) {
	c.Set("user_id", userID)
	c.Set("user_type", "c_end")
}

func decodeResponseMap(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

func seedHandlerStation(t *testing.T, db *gorm.DB, name string) *model.ServiceStation {
	t.Helper()

	station := &model.ServiceStation{
		Name:   name,
		Code:   name,
		Status: "active",
	}
	if err := db.Create(station).Error; err != nil {
		t.Fatalf("failed to create station: %v", err)
	}
	return station
}

func seedHandlerZone(t *testing.T, db *gorm.DB, stationID int64, name, points string) *model.ServiceZone {
	t.Helper()

	zone := &model.ServiceZone{
		StationID: stationID,
		Name:      name,
		Points:    points,
		Priority:  1,
		Status:    "active",
	}
	if err := db.Create(zone).Error; err != nil {
		t.Fatalf("failed to create zone: %v", err)
	}
	return zone
}

func seedHandlerRequest(t *testing.T, db *gorm.DB, userID, stationID int64, requestNo string) *model.ServiceRequest {
	t.Helper()

	req := &model.ServiceRequest{
		RequestNo:   requestNo,
		UserID:      userID,
		ServiceType: "meal",
		Status:      "dispatched",
		Address:     "测试地址",
		StationID:   stationID,
	}
	if err := db.Create(req).Error; err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	return req
}

func seedHandlerUserAndProfile(t *testing.T, db *gorm.DB, userID int64, phone, name, customerType, address string) {
	t.Helper()

	user := &model.User{
		ID:           userID,
		Phone:        phone,
		PasswordHash: "hash",
		Name:         name,
		Status:       "active",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	profile := &model.CustomerProfile{
		UserID:           userID,
		CustomerType:     customerType,
		Address:          address,
		EmergencyContact: `{}`,
	}
	if err := db.Create(profile).Error; err != nil {
		t.Fatalf("failed to create profile: %v", err)
	}
}
