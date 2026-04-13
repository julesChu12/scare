package handler

import "time"

// 通用响应结构
type APIResponse struct {
	Msg  string      `json:"msg" example:"ok"`
	Data interface{} `json:"data"`
}

// 分页响应
type PageResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total" example:"100"`
}

// ========== 用户模块 ==========

// UserResponse 用户响应
type UserResponse struct {
	ID              int64     `json:"id" example:"1"`
	Phone           string    `json:"phone" example:"13800000001"`
	Name            string    `json:"name" example:"张三"`
	Email           string    `json:"email" example:"zhangsan@example.com"`
	Avatar          string    `json:"avatar" example:"https://example.com/avatar.jpg"`
	Gender          string    `json:"gender" example:"male"`
	BirthDate       string    `json:"birth_date" example:"1990-01-01"`
	Age             *int      `json:"age,omitempty" example:"36"`
	IDCardMasked    string    `json:"id_card_masked,omitempty" example:"1101**********1234"`
	IDCardHash      string    `json:"id_card_hash,omitempty" example:"2f0ab5..."`
	IDCardToken     string    `json:"id_card_token,omitempty" example:"eyJ1aWQiOjF9.xxx"`
	StationID       int64     `json:"station_id" example:"1"`
	Status          string    `json:"status" example:"active"`
	PrimaryIdentity string    `json:"primary_identity" example:"staff"`
	BEndIdentities  []string  `json:"b_end_identities" example:"staff,station_manager"`
	CEndIdentities  []string  `json:"c_end_identities" example:"elderly"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Items []UserResponse `json:"items"`
	Total int64          `json:"total" example:"100"`
}

// ========== 服务请求模块 ==========

// RequestResponse 服务请求响应
type RequestResponse struct {
	ID                int64     `json:"id" example:"1"`
	RequestNo         string    `json:"request_no" example:"REQ20260209001"`
	UserID            int64     `json:"user_id" example:"8"`
	ServiceType       string    `json:"service_type" example:"meal"`
	Status            string    `json:"status" example:"dispatched"`
	Description       string    `json:"description" example:"需要午餐送餐"`
	ContactName       string    `json:"contact_name" example:"张大爷"`
	ContactPhone      string    `json:"contact_phone" example:"13800000008"`
	Address           string    `json:"address" example:"北京市朝阳区幸福小区1号楼"`
	SourceStationID   int64     `json:"source_station_id" example:"2"`
	StationID         int64     `json:"station_id" example:"1"`
	StationName       string    `json:"station_name" example:"朝阳区幸福街道服务站"`
	SourceStationName string    `json:"source_station_name" example:"霍营街道第一服务站"`
	DispatchBasis     string    `json:"dispatch_basis" example:"service_geofence"`
	NeedsManualVerify bool      `json:"needs_manual_verify" example:"false"`
	AppointmentTime   time.Time `json:"appointment_time"`
	Urgency           string    `json:"urgency" example:"normal"`
	RejectReason      string    `json:"reject_reason"`
	Images            string    `json:"images"`
	Rating            int64     `json:"rating" example:"5"`
	Feedback          string    `json:"feedback"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// RequestListResponse 服务请求列表响应
type RequestListResponse struct {
	Items    []RequestResponse `json:"items"`
	Total    int64             `json:"total" example:"100"`
	Page     int               `json:"page" example:"1"`
	PageSize int               `json:"page_size" example:"10"`
}

// ========== 任务模块 ==========

// TaskResponse 任务响应
type TaskResponse struct {
	ID          int64            `json:"id" example:"1"`
	RequestID   int64            `json:"request_id" example:"1"`
	StationID   int64            `json:"station_id" example:"1"`
	StaffID     int64            `json:"staff_id" example:"4"`
	Status      string           `json:"status" example:"claimed"`
	ClaimedAt   time.Time        `json:"claimed_at"`
	CompletedAt time.Time        `json:"completed_at"`
	Rating      int64            `json:"rating" example:"5"`
	Feedback    string           `json:"feedback"`
	StaffNotes  string           `json:"staff_notes"`
	Images      string           `json:"images"`
	Request     *RequestResponse `json:"request,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Items    []TaskResponse `json:"items"`
	Total    int64          `json:"total" example:"100"`
	Page     int            `json:"page" example:"1"`
	PageSize int            `json:"page_size" example:"10"`
}

// ========== 站点模块 ==========

// StationResponse 站点响应
type StationResponse struct {
	ID          int64     `json:"id" example:"1"`
	Name        string    `json:"name" example:"朝阳区幸福街道服务站"`
	Code        string    `json:"code" example:"ST001"`
	Address     string    `json:"address" example:"北京市朝阳区幸福街道1号"`
	Phone       string    `json:"phone" example:"010-12345678"`
	Latitude    float64   `json:"latitude" example:"39.9042"`
	Longitude   float64   `json:"longitude" example:"116.4074"`
	ServiceArea string    `json:"service_area" example:"幸福街道辖区"`
	Capacity    int64     `json:"capacity" example:"10"`
	WorkHours   string    `json:"work_hours" example:"08:00-18:00"`
	Status      string    `json:"status" example:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StationListResponse 站点列表响应
type StationListResponse struct {
	Items []StationResponse `json:"items"`
	Total int64             `json:"total" example:"10"`
}

// ========== 权限模块 ==========

// PermissionNode 权限树节点
type PermissionNode struct {
	ID       string           `json:"id" example:"service:request"`
	Code     string           `json:"code" example:"service:request"`
	Label    string           `json:"label" example:"服务请求"`
	Type     string           `json:"type" example:"menu"`
	APIPath  string           `json:"api_path,omitempty" example:"/api/v1/b/requests"`
	Method   string           `json:"method,omitempty" example:"GET"`
	IsPublic bool             `json:"is_public,omitempty" example:"false"`
	Disabled bool             `json:"disabled,omitempty" example:"false"`
	Children []PermissionNode `json:"children,omitempty"`
}

// PermissionTreeResponse 权限树响应
type PermissionTreeResponse struct {
	Tree []PermissionNode `json:"tree"`
}

// ========== 认证模块 ==========

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string       `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	User         UserResponse `json:"user"`
}

// ========== 上传模块 ==========

// UploadResponse 上传响应
type UploadResponse struct {
	URL string `json:"url" example:"http://localhost:8080/static/uploads/xxx.jpg"`
	Key string `json:"key" example:"uploads/xxx.jpg"`
}

// ========== 统计模块 ==========

// DashboardStatsResponse 仪表盘统计响应
type DashboardStatsResponse struct {
	TotalUsers       int64 `json:"total_users" example:"100"`
	TotalStations    int64 `json:"total_stations" example:"5"`
	TotalRequests    int64 `json:"total_requests" example:"500"`
	PendingRequests  int64 `json:"pending_requests" example:"20"`
	TodayNewRequests int64 `json:"today_new_requests" example:"10"`
	TodayCompleted   int64 `json:"today_completed" example:"8"`
}
