package consts

// B端身份类型（参与 RBAC 权限验证）
const (
	IdentityAdmin          = "admin"
	IdentityStationManager = "station_manager"
	IdentityStaff          = "staff"
)

// C端身份类型（不参与 RBAC 权限验证）
const (
	IdentityElderly  = "elderly"
	IdentityFamily   = "family"
	IdentityPregnant = "pregnant"
	IdentityDisabled = "disabled"
	IdentityChild    = "child"
)

// 兼容旧常量（逐步弃用）
const (
	RoleElderly        = IdentityElderly
	RoleFamily         = IdentityFamily
	RoleStaff          = IdentityStaff
	RoleStationManager = IdentityStationManager
	RoleAdmin          = IdentityAdmin
)

// BEndIdentityTypes B端身份类型列表
var BEndIdentityTypes = []string{
	IdentityAdmin,
	IdentityStationManager,
	IdentityStaff,
}

// CEndIdentityTypes C端身份类型列表
var CEndIdentityTypes = []string{
	IdentityElderly,
	IdentityFamily,
	IdentityPregnant,
	IdentityDisabled,
	IdentityChild,
}

// IsBEndIdentity 判断是否为 B端身份
func IsBEndIdentity(identityType string) bool {
	for _, t := range BEndIdentityTypes {
		if t == identityType {
			return true
		}
	}
	return false
}

// IsCEndIdentity 判断是否为 C端身份
func IsCEndIdentity(identityType string) bool {
	for _, t := range CEndIdentityTypes {
		if t == identityType {
			return true
		}
	}
	return false
}
