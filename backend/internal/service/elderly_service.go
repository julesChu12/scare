package service

import (
	"errors"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrElderlyNotFound  = errors.New("elderly profile not found")
	ErrPhoneDuplicate   = errors.New("phone number already registered")
	ErrInvalidStationID = errors.New("invalid station id")
)

// ElderlyService B端老年人档案服务
type ElderlyService struct {
	db           *gorm.DB
	userRepo     *repository.UserRepository
	identityRepo *repository.UserIdentityRepository
	customerRepo *repository.CustomerRepository
	requestRepo  *repository.RequestRepository
	taskRepo     *repository.TaskRepository
}

func NewElderlyService(
	db *gorm.DB,
	userRepo *repository.UserRepository,
	identityRepo *repository.UserIdentityRepository,
	customerRepo *repository.CustomerRepository,
	requestRepo *repository.RequestRepository,
	taskRepo *repository.TaskRepository,
) *ElderlyService {
	return &ElderlyService{
		db:           db,
		userRepo:     userRepo,
		identityRepo: identityRepo,
		customerRepo: customerRepo,
		requestRepo:  requestRepo,
		taskRepo:     taskRepo,
	}
}

// ElderlyInput 创建/更新老人档案的输入
type ElderlyInput struct {
	Name            string    `json:"name"`
	Phone           string    `json:"phone"`
	Gender          string    `json:"gender"`
	BirthDate       time.Time `json:"birth_date"`
	IDCard          string    `json:"id_card"`
	Address         string    `json:"address"`
	StationID       int64     `json:"station_id"`
	HealthStatus    string    `json:"health_status"`
	DisabilityLevel string    `json:"disability_level"`
	MedicalHistory  string    `json:"medical_history"`
	SpecialNeeds    string    `json:"special_needs"`
}

// ElderlyInfo 老人档案完整信息
type ElderlyInfo struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Phone           string    `json:"phone"`
	Gender          string    `json:"gender"`
	BirthDate       string    `json:"birth_date"`
	IDCard          string    `json:"id_card"`
	Address         string    `json:"address"`
	StationID       int64     `json:"station_id"`
	StationName     string    `json:"station_name"`
	HealthStatus    string    `json:"health_status"`
	DisabilityLevel string    `json:"disability_level"`
	MedicalHistory  string    `json:"medical_history"`
	SpecialNeeds    string    `json:"special_needs"`
	CustomerType    string    `json:"customer_type"`
	CreatedAt       time.Time `json:"created_at"`
}

// ElderlyFilter 老人档案列表筛选条件
type ElderlyFilter struct {
	Keyword      string
	StationID    int64
	HealthStatus string
	Page         int
	PageSize     int
}

// List 获取老人档案列表
func (s *ElderlyService) List(filter ElderlyFilter) ([]*ElderlyInfo, int64, error) {
	db := s.db.Table("users u").
		Select(`u.id, u.name, u.phone, u.gender, u.birth_date, u.id_card, u.station_id, u.created_at,
			cp.address, cp.health_status, cp.disability_level, cp.medical_history, cp.special_needs, cp.customer_type,
			ss.name as station_name`).
		Joins("INNER JOIN user_identities ui ON ui.user_id = u.id AND ui.identity_type = ? AND ui.deleted_at IS NULL", consts.IdentityElderly).
		Joins("LEFT JOIN customer_profiles cp ON cp.user_id = u.id AND cp.deleted_at IS NULL").
		Joins("LEFT JOIN service_stations ss ON u.station_id = ss.id AND ss.deleted_at IS NULL").
		Where("u.deleted_at IS NULL")

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		db = db.Where("(u.name LIKE ? OR u.phone LIKE ?)", keyword, keyword)
	}
	if filter.StationID > 0 {
		db = db.Where("u.station_id = ?", filter.StationID)
	}
	if filter.HealthStatus != "" {
		db = db.Where("cp.health_status = ?", filter.HealthStatus)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	var items []*ElderlyInfo
	if err := db.Order("u.id DESC").Offset(offset).Limit(filter.PageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	// 格式化 birth_date
	for _, item := range items {
		if item.BirthDate == "0001-01-01" || item.BirthDate == "0000-00-00" {
			item.BirthDate = ""
		}
	}

	return items, total, nil
}

// GetByID 获取老人档案详情
func (s *ElderlyService) GetByID(userID int64) (*ElderlyInfo, error) {
	var info ElderlyInfo
	err := s.db.Table("users u").
		Select(`u.id, u.name, u.phone, u.gender, u.birth_date, u.id_card, u.station_id, u.created_at,
			cp.address, cp.health_status, cp.disability_level, cp.medical_history, cp.special_needs, cp.customer_type,
			ss.name as station_name`).
		Joins("INNER JOIN user_identities ui ON ui.user_id = u.id AND ui.identity_type = ? AND ui.deleted_at IS NULL", consts.IdentityElderly).
		Joins("LEFT JOIN customer_profiles cp ON cp.user_id = u.id AND cp.deleted_at IS NULL").
		Joins("LEFT JOIN service_stations ss ON u.station_id = ss.id AND ss.deleted_at IS NULL").
		Where("u.id = ? AND u.deleted_at IS NULL", userID).
		First(&info).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrElderlyNotFound
		}
		return nil, err
	}

	if info.BirthDate == "0001-01-01" || info.BirthDate == "0000-00-00" {
		info.BirthDate = ""
	}

	return &info, nil
}

// Create 创建老人档案
func (s *ElderlyService) Create(input ElderlyInput) (*ElderlyInfo, error) {
	if input.Phone == "" || input.Name == "" {
		return nil, ErrInvalidUser
	}

	// 检查手机号唯一性
	existing, _ := s.userRepo.GetByPhone(input.Phone)
	if existing != nil {
		return nil, ErrPhoneDuplicate
	}

	// 默认密码
	hash, err := HashPassword("Elderly@123")
	if err != nil {
		return nil, err
	}

	var userID int64

	// 事务内创建 user + identity + profile
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建 user
		user := &model.User{
			Phone:        input.Phone,
			PasswordHash: hash,
			Name:         input.Name,
			Gender:       input.Gender,
			BirthDate:    input.BirthDate,
			IDCard:       input.IDCard,
			StationID:    input.StationID,
			Status:       "active",
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		userID = user.ID

		// 2. 创建 user_identity
		identity := &model.UserIdentity{
			UserID:       user.ID,
			IdentityType: consts.IdentityElderly,
			IsPrimary:    true,
			Status:       "active",
		}
		if err := tx.Create(identity).Error; err != nil {
			return err
		}

		// 3. 创建 customer_profile
		healthStatus := input.HealthStatus
		if healthStatus == "" {
			healthStatus = "good"
		}
		profile := &model.CustomerProfile{
			UserID:          user.ID,
			Address:         input.Address,
			CustomerType:    consts.IdentityElderly,
			HealthStatus:    healthStatus,
			DisabilityLevel: input.DisabilityLevel,
			MedicalHistory:  input.MedicalHistory,
			SpecialNeeds:    input.SpecialNeeds,
		}
		if err := tx.Create(profile).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetByID(userID)
}

// Update 更新老人档案
func (s *ElderlyService) Update(userID int64, input ElderlyInput) (*ElderlyInfo, error) {
	// 先确认老人存在
	_, err := s.GetByID(userID)
	if err != nil {
		return nil, err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 更新 users 表字段
		userUpdates := map[string]interface{}{}
		if input.Name != "" {
			userUpdates["name"] = input.Name
		}
		if input.Gender != "" {
			userUpdates["gender"] = input.Gender
		}
		if !input.BirthDate.IsZero() {
			userUpdates["birth_date"] = input.BirthDate
		}
		if input.IDCard != "" {
			userUpdates["id_card"] = input.IDCard
		}
		if input.StationID > 0 {
			userUpdates["station_id"] = input.StationID
		}
		if len(userUpdates) > 0 {
			if err := tx.Table("users").Where("id = ?", userID).Updates(userUpdates).Error; err != nil {
				return err
			}
		}

		// 更新 customer_profiles 表字段
		profileUpdates := map[string]interface{}{}
		if input.Address != "" {
			profileUpdates["address"] = input.Address
		}
		if input.HealthStatus != "" {
			profileUpdates["health_status"] = input.HealthStatus
		}
		// 以下字段允许设为空字符串，所以始终更新
		profileUpdates["disability_level"] = input.DisabilityLevel
		profileUpdates["medical_history"] = input.MedicalHistory
		profileUpdates["special_needs"] = input.SpecialNeeds

		if err := tx.Table("customer_profiles").Where("user_id = ?", userID).Updates(profileUpdates).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.GetByID(userID)
}

// ServiceRecordInfo 服务记录信息
type ServiceRecordInfo struct {
	RequestID   int64     `json:"request_id"`
	RequestNo   string    `json:"request_no"`
	ServiceType string    `json:"service_type"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Rating      int64     `json:"rating"`
	Feedback    string    `json:"feedback"`
	CreatedAt   time.Time `json:"created_at"`
	// 任务信息（如有）
	TaskID      *int64     `json:"task_id"`
	StaffName   *string    `json:"staff_name"`
	ClaimedAt   *time.Time `json:"claimed_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// GetServiceRecords 获取老人的服务记录
func (s *ElderlyService) GetServiceRecords(userID int64, page, pageSize int) ([]*ServiceRecordInfo, int64, error) {
	db := s.db.Table("service_requests sr").
		Select(`sr.id as request_id, sr.request_no, sr.service_type, sr.status,
			sr.description, sr.address, sr.rating, sr.feedback, sr.created_at,
			t.id as task_id, u.name as staff_name, t.claimed_at, t.completed_at`).
		Joins("LEFT JOIN tasks t ON t.request_id = sr.id AND t.deleted_at IS NULL").
		Joins("LEFT JOIN users u ON u.id = t.staff_id AND u.deleted_at IS NULL").
		Where("sr.user_id = ? AND sr.deleted_at IS NULL", userID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var items []*ServiceRecordInfo
	if err := db.Order("sr.created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
