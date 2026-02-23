package handler

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (int64, bool) {
	value, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	id, ok := value.(int64)
	return id, ok
}

func GetUserType(c *gin.Context) (string, bool) {
	value, ok := c.Get("user_type")
	if !ok {
		return "", false
	}
	userType, ok := value.(string)
	return userType, ok
}

func GetStationID(c *gin.Context) (int64, bool) {
	value, ok := c.Get("station_id")
	if !ok {
		return 0, false
	}
	id, ok := value.(int64)
	return id, ok
}

// GetUserIdentities 获取用户身份列表
func GetUserIdentities(c *gin.Context) []string {
	value, ok := c.Get("user_identities")
	if !ok {
		return nil
	}
	identities, ok := value.([]string)
	if !ok {
		return nil
	}
	return identities
}

// GetUserRoles 获取用户角色列表（兼容旧代码，返回身份列表）
// Deprecated: use GetUserIdentities instead
func GetUserRoles(c *gin.Context) []string {
	return GetUserIdentities(c)
}
