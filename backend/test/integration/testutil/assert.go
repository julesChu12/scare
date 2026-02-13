//go:build integration

package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// APIResponse 通用 API 响应结构
type APIResponse struct {
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// AssertOK 验证 200 + msg:"ok" + data 非 nil
func AssertOK(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "ok", resp.Msg)
	assert.NotNil(t, resp.Data)

	var data map[string]interface{}
	_ = json.Unmarshal(resp.Data, &data)
	return data
}

// AssertMsg 验证 200 + 指定 msg
func AssertMsg(t *testing.T, w *httptest.ResponseRecorder, msg string) {
	t.Helper()
	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, msg, resp.Msg)
}

// AssertPageResponse 验证分页结构
func AssertPageResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	data := AssertOK(t, w)
	assert.Contains(t, data, "items")
	assert.Contains(t, data, "page")
	assert.Contains(t, data, "page_size")
	assert.Contains(t, data, "total")
	return data
}

// AssertError 验证错误响应
func AssertError(t *testing.T, w *httptest.ResponseRecorder, code int, msgContains string) {
	t.Helper()
	assert.Equal(t, code, w.Code)
	var resp APIResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp.Msg, msgContains)
}

// ParseData 解析 response.data 为 map
func ParseData(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var resp APIResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	var data map[string]interface{}
	_ = json.Unmarshal(resp.Data, &data)
	return data
}

// ParseDataSlice 解析 response.data 为 slice
func ParseDataSlice(t *testing.T, w *httptest.ResponseRecorder) []interface{} {
	t.Helper()
	var resp APIResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	var data []interface{}
	_ = json.Unmarshal(resp.Data, &data)
	return data
}

// DoRequest 发送 HTTP 请求并返回 recorder
func DoRequest(engine http.Handler, method, path, token string, body ...string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var req *http.Request
	if len(body) > 0 && body[0] != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body[0]))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	engine.ServeHTTP(w, req)
	return w
}
