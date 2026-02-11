package consts

// 服务类型常量定义
const (
	ServiceTypeMeal       = "meal"       // 送餐服务
	ServiceTypeMedical    = "medical"    // 就医陪护
	ServiceTypeCare       = "care"       // 日常照护
	ServiceTypeRepair     = "repair"     // 居家维修
	ServiceTypeCleaning   = "cleaning"   // 家政保洁
	ServiceTypeCompany    = "company"    // 陪伴聊天
	ServiceTypeEmergency  = "emergency"  // 紧急救助
	ServiceTypeShopping   = "shopping"   // 代买代购
	ServiceTypeTransport  = "transport"  // 出行接送
	ServiceTypeRehab      = "rehab"      // 康复理疗
	ServiceTypePsychology = "psychology" // 心理慰藉
	ServiceTypeLegalAid   = "legal_aid"  // 法律援助
	ServiceTypeOther      = "other"      // 其他服务
)

// ServiceTypes 所有服务类型列表
var ServiceTypes = []string{
	ServiceTypeMeal,
	ServiceTypeMedical,
	ServiceTypeCare,
	ServiceTypeRepair,
	ServiceTypeCleaning,
	ServiceTypeCompany,
	ServiceTypeEmergency,
	ServiceTypeShopping,
	ServiceTypeTransport,
	ServiceTypeRehab,
	ServiceTypePsychology,
	ServiceTypeLegalAid,
	ServiceTypeOther,
}

// ServiceTypeNames 服务类型显示名称映射
var ServiceTypeNames = map[string]string{
	ServiceTypeMeal:       "送餐服务",
	ServiceTypeMedical:    "就医陪护",
	ServiceTypeCare:       "日常照护",
	ServiceTypeRepair:     "居家维修",
	ServiceTypeCleaning:   "家政保洁",
	ServiceTypeCompany:    "陪伴聊天",
	ServiceTypeEmergency:  "紧急救助",
	ServiceTypeShopping:   "代买代购",
	ServiceTypeTransport:  "出行接送",
	ServiceTypeRehab:      "康复理疗",
	ServiceTypePsychology: "心理慰藉",
	ServiceTypeLegalAid:   "法律援助",
	ServiceTypeOther:      "其他服务",
}

// IsValidServiceType 验证服务类型是否合法
func IsValidServiceType(serviceType string) bool {
	for _, t := range ServiceTypes {
		if t == serviceType {
			return true
		}
	}
	return false
}

// GetServiceTypeName 获取服务类型显示名称
func GetServiceTypeName(serviceType string) string {
	if name, ok := ServiceTypeNames[serviceType]; ok {
		return name
	}
	return "未知服务"
}
