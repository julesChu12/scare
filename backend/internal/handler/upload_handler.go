package handler

import (
	"net/http"

	"community-elderly-care-platform/internal/service"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service *service.StorageService
}

func NewUploadHandler(service *service.StorageService) *UploadHandler {
	return &UploadHandler{service: service}
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
// @Failure      500  {object} APIResponse "上传失败"
// @Router       /c/upload [post]
// @Router       /b/upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		RespondError(c, http.StatusBadRequest, "file required")
		return
	}

	userType, _ := GetUserType(c)
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

	Respond(c, http.StatusOK, "ok", gin.H{
		"url": url,
		"key": key,
	})
}
