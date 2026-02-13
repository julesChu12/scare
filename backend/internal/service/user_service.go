package service

import (
	"errors"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/internal/repository"
)

var ErrInvalidUser = errors.New("invalid user")

// 允许的身份类型
var allowedIdentityTypes = map[string]struct{}{
	consts.IdentityElderly:        {},
	consts.IdentityFamily:         {},
	consts.IdentityPregnant:       {},
	consts.IdentityDisabled:       {},
	consts.IdentityChild:          {},
	consts.IdentityStaff:          {},
	consts.IdentityStationManager: {},
	consts.IdentityAdmin:          {},
}

type UserService struct {
	repo         *repository.UserRepository
	identityRepo *repository.UserIdentityRepository
}

// UserInput 创建/更新用户的输入参数
type UserInput struct {
	ID           int64     `json:"id"`             // 用户ID（更新时必填）
	Phone        string    `json:"phone"`          // 手机号
	Password     string    `json:"password"`       // 密码
	Name         string    `json:"name"`           // 姓名
	Email        string    `json:"email"`          // 邮箱
	Avatar       string    `json:"avatar"`         // 头像
	Gender       string    `json:"gender"`         // 性别
	BirthDate    time.Time `json:"birth_date"`     // 出生日期
	IDCard       string    `json:"id_card"`        // 身份证号（密文或明文）
	IDCardHMAC   string    `json:"id_card_hmac"`   // 身份证号HMAC摘要
	IDCardMasked string    `json:"id_card_masked"` // 身份证号脱敏值
	IdentityType string    `json:"identity_type"`  // 身份类型
	StationID    int64     `json:"station_id"`     // 站点ID
	Status       string    `json:"status"`         // 状态
}

// UserFilter 用户筛选条件（Service 层）
type UserFilter struct {
	Role      string // 按身份类型筛选
	Status    string // 按状态筛选
	StationID int64  // 按站点筛选
	Keyword   string // 关键词搜索
}

// UserWithIdentities 用户及其身份信息
type UserWithIdentities struct {
	*model.User
	Identities      []*model.UserIdentity // 所有身份
	PrimaryIdentity *model.UserIdentity   // 主身份
	BEndIdentities  []string              // B端身份列表
	CEndIdentities  []string              // C端身份列表
}

func NewUserService(repo *repository.UserRepository, identityRepo *repository.UserIdentityRepository) *UserService {
	return &UserService{
		repo:         repo,
		identityRepo: identityRepo,
	}
}

func (s *UserService) Create(input UserInput) (*UserWithIdentities, error) {
	if input.Phone == "" || input.Password == "" || input.IdentityType == "" {
		return nil, ErrInvalidUser
	}
	if _, ok := allowedIdentityTypes[input.IdentityType]; !ok {
		return nil, ErrInvalidUser
	}

	hash, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Phone:        input.Phone,
		PasswordHash: hash,
		Name:         input.Name,
		Email:        input.Email,
		StationID:    input.StationID,
		Status:       input.Status,
	}
	if user.Status == "" {
		user.Status = "active"
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// 创建用户身份
	identity := &model.UserIdentity{
		UserID:       user.ID,
		IdentityType: input.IdentityType,
		IsPrimary:    true,
		Status:       "active",
	}
	// B端身份需要站点
	if consts.IsBEndIdentity(input.IdentityType) {
		identity.StationID = input.StationID
	}

	if err := s.identityRepo.Create(identity); err != nil {
		return nil, err
	}

	return s.buildUserWithIdentities(user)
}

func (s *UserService) Update(input UserInput) (*UserWithIdentities, error) {
	if input.ID == 0 {
		return nil, ErrInvalidUser
	}
	user, err := s.repo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Avatar != "" {
		user.Avatar = input.Avatar
	}
	if input.Gender != "" {
		user.Gender = input.Gender
	}
	if !input.BirthDate.IsZero() {
		user.BirthDate = input.BirthDate
	}
	if input.IDCard != "" {
		user.IDCard = input.IDCard
		user.IDCardHmac = input.IDCardHMAC
		user.IDCardMasked = input.IDCardMasked
	}
	if input.Status != "" {
		user.Status = input.Status
	}
	if input.StationID != 0 {
		user.StationID = input.StationID
	}
	if input.Password != "" {
		hash, err := HashPassword(input.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.buildUserWithIdentities(user)
}

func (s *UserService) List(page, pageSize int) ([]*UserWithIdentities, int64, error) {
	return s.ListWithFilter(page, pageSize, UserFilter{})
}

// ListWithFilter 获取用户列表（带筛选条件）
func (s *UserService) ListWithFilter(page, pageSize int, filter UserFilter) ([]*UserWithIdentities, int64, error) {
	offset := (page - 1) * pageSize

	// 转换为 repository 层的 filter
	repoFilter := repository.UserFilter{
		Role:      filter.Role,
		Status:    filter.Status,
		StationID: filter.StationID,
		Keyword:   filter.Keyword,
	}

	users, total, err := s.repo.ListWithFilter(offset, pageSize, repoFilter)
	if err != nil {
		return nil, 0, err
	}

	// 组装用户身份信息
	result := make([]*UserWithIdentities, 0, len(users))
	for _, user := range users {
		uwi, err := s.buildUserWithIdentities(user)
		if err != nil {
			continue
		}
		result = append(result, uwi)
	}

	return result, total, nil
}

func (s *UserService) GetByID(id int64) (*UserWithIdentities, error) {
	if id == 0 {
		return nil, ErrInvalidUser
	}
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.buildUserWithIdentities(user)
}

// GetByPhone 根据手机号获取用户及身份
func (s *UserService) GetByPhone(phone string) (*UserWithIdentities, error) {
	user, err := s.repo.GetByPhone(phone)
	if err != nil {
		return nil, err
	}
	return s.buildUserWithIdentities(user)
}

// buildUserWithIdentities 构建用户身份信息
func (s *UserService) buildUserWithIdentities(user *model.User) (*UserWithIdentities, error) {
	identities, err := s.identityRepo.GetActiveByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	uwi := &UserWithIdentities{
		User:           user,
		Identities:     identities,
		BEndIdentities: FilterBEndIdentities(identities),
		CEndIdentities: FilterCEndIdentities(identities),
	}

	// 获取主身份
	for _, identity := range identities {
		if identity.IsPrimary {
			uwi.PrimaryIdentity = identity
			break
		}
	}

	return uwi, nil
}
