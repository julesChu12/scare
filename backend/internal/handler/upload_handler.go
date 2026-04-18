package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"community-elderly-care-platform/internal/service"
	"community-elderly-care-platform/pkg/redis"

	"github.com/gin-gonic/gin"
)

var ErrUploadRateLimit = errors.New("上传过于频繁，请稍后再试")

type UploadHandler struct {
	service *service.StorageService
	rdb     *redis.Client
}

// NewUploadHandler 创建 UploadHandler
func NewUploadHandler(service *service.StorageService, rdb *redis.Client) *UploadHandler {
	return &UploadHandler{service: service, rdb: rdb}
}

// checkRateLimit 检查上传频率限制（C端：IP维度，每分钟20次；B端：用户维度，每分钟50次）
func (h *UploadHandler) checkRateLimit(_ *gin.Context, userID int64, userType string) error {
	ctx := context.Background()
	key := fmt.Sprintf("upload_rate:%s:%d", userType, userID)

	if h.rdb == nil {
		return nil // 无 Redis 时跳过限流
	}

	// 检查是否已有限制键
	exists, err := h.rdb.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return ErrUploadRateLimit
	}

	return nil
}

// incrementRateLimit 增加上传频率计数
func (h *UploadHandler) incrementRateLimit(_ *gin.Context, userID int64, userType string) error {
	if h.rdb == nil {
		return nil
	}
	ctx := context.Background()
	key := fmt.Sprintf("upload_rate:%s:%d", userType, userID)
	return h.rdb.SetEx(ctx, key, "1", time.Minute).Err()
}

// Upload 上传文件
// @Summary      上传文件
// @Description  上传文件到存储服务，B端和C端有不同的限制
// @Tags         c_upload,b_upload
// @Accept       multipart/form-data
// @Produce      json
// @Security     Bearer
// @Param        file formData file true "文件"
// @Param        module query string false "模块名称"
// @Success      200  {object} APIResponse "上传成功，返回url和key"
// @Failure      400  {object} APIResponse "文件缺失或不符合要求"
// @Failure      401  {object} APIResponse "未认证"
// @Failure      429  {object} APIResponse "上传过于频繁"
// @Failure      500  {object} APIResponse "上传失败"
// @Router       /c/upload [post]
// @Router       /b/upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		RespondError(c, http.StatusBadRequest, "file required")
		return
	}

	userID, _ := GetUserID(c)
	userType, _ := GetUserType(c)

	// 频率限制：未登录用户用 IP 作为标识
	identifier := userID
	if userID == 0 {
		identifier = int64(c.Request.RemoteAddr[0] << 24)
	}

	if err := h.checkRateLimit(c, identifier, userType); err != nil {
		if errors.Is(err, ErrUploadRateLimit) {
			RespondError(c, http.StatusTooManyRequests, err.Error())
			return
		}
		// Redis 错误不阻止上传，只记录日志
	}

	userType, _ = GetUserType(c)
	module := c.Query("module")

	// C端：仅允许图片，限制大小 5MB
	if userType == "c_end" {
		if file.Size > 5*1024*1024 {
			RespondError(c, http.StatusBadRequest, "file too large (max 5MB)")
			return
		}
		contentType := file.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
			RespondError(c, http.StatusBadRequest, "only image files allowed")
			return
		}
		if module == "" {
			module = "c_end"
		}
	} else {
		// B端：允许多种文件，限制大小 50MB
		if file.Size > 50*1024*1024 {
			RespondError(c, http.StatusBadRequest, "file too large (max 50MB)")
			return
		}
		if module == "" {
			module = "b_end"
		}
	}

	url, key, err := h.service.Upload(c.Request.Context(), module, file)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "upload failed")
		return
	}

	// 增加频率计数
	_ = h.incrementRateLimit(c, identifier, userType)

	// 返回相对路径给前端，前端存相对路径到数据库，API 返回时再动态拼接完整 URL
	Respond(c, http.StatusOK, "ok", gin.H{
		"url": url,
		"key": key,
	})
}
