//go:build integration

package integration

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestUpload(t *testing.T) {
	env := testutil.Setup(t)

	t.Run("B端文件上传", func(t *testing.T) {
		// 构造 multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.txt")
		part.Write([]byte("test file content"))
		writer.Close()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/b/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+testutil.AdminToken())
		env.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("C端文件上传", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		header := textproto.MIMEHeader{}
		header.Set("Content-Disposition", `form-data; name="file"; filename="photo.jpg"`)
		header.Set("Content-Type", "image/jpeg")
		part, _ := writer.CreatePart(header)
		part.Write([]byte("fake image data"))
		writer.Close()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/c/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+testutil.CEndElderlyToken())
		env.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("未认证上传被拒绝", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.txt")
		part.Write([]byte("test"))
		writer.Close()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/b/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		env.Engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
