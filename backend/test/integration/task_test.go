//go:build integration

package integration

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"community-elderly-care-platform/internal/consts"
	"community-elderly-care-platform/internal/dao/model"
	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	env := testutil.Setup(t)
	adminToken := testutil.AdminToken()
	stationManagerToken := testutil.StationManagerToken()
	staffToken := testutil.StaffToken()

	t.Run("B端_任务列表_Admin", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks", adminToken)
		testutil.AssertOK(t, w)
	})

	t.Run("B端_任务池_Admin", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks/pool", adminToken)
		testutil.AssertOK(t, w)
	})

	t.Run("B端_我的任务_Staff", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks/my", staffToken)
		testutil.AssertOK(t, w)
	})

	t.Run("B端_任务详情_Admin", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks/1", adminToken)
		// 种子数据中无任务记录，404 是预期行为
		assert.Contains(t,
			[]int{http.StatusOK, http.StatusNotFound}, w.Code)
	})

	t.Run("B端_任务详情_Staff_可访问本站点任务", func(t *testing.T) {
		now := time.Now()
		req := &model.ServiceRequest{
			RequestNo:         "REQ-DETAIL-001",
			UserID:            10,
			ServiceType:       "meal",
			Status:            consts.RequestStatusClaimed,
			ContactName:       "李大爷",
			ContactPhone:      "13900000002",
			Address:           "测试地址详情",
			StationID:         1,
			CreatedAt:         now,
			UpdatedAt:         now,
			SubmitLocationLat: 39.9,
			SubmitLocationLng: 116.4,
		}
		assert.NoError(t, env.DB.DB.Create(req).Error)

		task := &model.TaskAssignment{
			RequestID: req.ID,
			StationID: 1,
			StaffID:   4,
			Status:    consts.TaskStatusClaimed,
			CreatedAt: now,
			UpdatedAt: now,
			ClaimedAt: now,
		}
		assert.NoError(t, env.DB.DB.Create(task).Error)

		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks/"+strconv.FormatInt(task.ID, 10), staffToken)
		data := testutil.AssertOK(t, w)
		assert.Equal(t, float64(task.ID), data["id"])
		request, ok := data["request"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "REQ-DETAIL-001", request["request_no"])
	})

	t.Run("Staff_无权访问任务池", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks/pool", staffToken)
		testutil.AssertError(t, w, http.StatusForbidden, "forbidden")
	})

	t.Run("B端_我的任务_支持需求编号筛选", func(t *testing.T) {
		now := time.Now()
		req1 := &model.ServiceRequest{
			RequestNo:         "REQ-FILTER-001",
			UserID:            10,
			ServiceType:       "meal",
			Status:            consts.RequestStatusClaimed,
			ContactName:       "张大爷",
			ContactPhone:      "13900000001",
			Address:           "测试地址1",
			StationID:         1,
			CreatedAt:         now,
			UpdatedAt:         now,
			SubmitLocationLat: 39.9,
			SubmitLocationLng: 116.4,
		}
		req2 := &model.ServiceRequest{
			RequestNo:         "REQ-OTHER-002",
			UserID:            10,
			ServiceType:       "medical",
			Status:            consts.RequestStatusClaimed,
			ContactName:       "张大爷",
			ContactPhone:      "13900000001",
			Address:           "测试地址2",
			StationID:         1,
			CreatedAt:         now,
			UpdatedAt:         now,
			SubmitLocationLat: 39.9,
			SubmitLocationLng: 116.4,
		}
		assert.NoError(t, env.DB.DB.Create(req1).Error)
		assert.NoError(t, env.DB.DB.Create(req2).Error)

		task1 := &model.TaskAssignment{
			RequestID: req1.ID,
			StationID: 1,
			StaffID:   4,
			Status:    consts.TaskStatusClaimed,
			CreatedAt: now,
			UpdatedAt: now,
			ClaimedAt: now,
		}
		task2 := &model.TaskAssignment{
			RequestID: req2.ID,
			StationID: 1,
			StaffID:   4,
			Status:    consts.TaskStatusClaimed,
			CreatedAt: now,
			UpdatedAt: now,
			ClaimedAt: now,
		}
		assert.NoError(t, env.DB.DB.Create(task1).Error)
		assert.NoError(t, env.DB.DB.Create(task2).Error)

		w := testutil.DoRequest(env.Engine, http.MethodGet,
			"/api/v1/b/tasks/my?page=1&page_size=10&request_no=REQ-FILTER-001", staffToken)
		data := testutil.AssertPageResponse(t, w)
		assert.Equal(t, float64(1), data["total"])

		items, ok := data["items"].([]interface{})
		assert.True(t, ok)
		if assert.Len(t, items, 1) {
			item, ok := items[0].(map[string]interface{})
			assert.True(t, ok)
			request, ok := item["request"].(map[string]interface{})
			assert.True(t, ok)
			assert.Equal(t, "REQ-FILTER-001", request["request_no"])
		}
	})

	t.Run("B端_完成任务_站点管理员可兜底完成本站任务", func(t *testing.T) {
		now := time.Now()
		req := &model.ServiceRequest{
			RequestNo:         "REQ-COMPLETE-SM-001",
			UserID:            10,
			ServiceType:       "meal",
			Status:            consts.RequestStatusClaimed,
			ContactName:       "张大爷",
			ContactPhone:      "13900000001",
			Address:           "测试地址-站长完成",
			StationID:         1,
			CreatedAt:         now,
			UpdatedAt:         now,
			SubmitLocationLat: 39.9,
			SubmitLocationLng: 116.4,
		}
		assert.NoError(t, env.DB.DB.Create(req).Error)

		task := &model.TaskAssignment{
			RequestID: req.ID,
			StationID: 1,
			StaffID:   4,
			Status:    consts.TaskStatusClaimed,
			CreatedAt: now,
			UpdatedAt: now,
			ClaimedAt: now,
		}
		assert.NoError(t, env.DB.DB.Create(task).Error)

		w := testutil.DoRequest(env.Engine, http.MethodPost,
			"/api/v1/b/tasks/"+strconv.FormatInt(task.ID, 10)+"/complete", stationManagerToken, `{"images":["proof-sm.jpg"]}`)
		testutil.AssertOK(t, w)

		var updatedTask model.TaskAssignment
		assert.NoError(t, env.DB.DB.First(&updatedTask, task.ID).Error)
		assert.Equal(t, consts.TaskStatusCompleted, updatedTask.Status)
		assert.NotZero(t, updatedTask.CompletedAt)

		var updatedReq model.ServiceRequest
		assert.NoError(t, env.DB.DB.First(&updatedReq, req.ID).Error)
		assert.Equal(t, consts.RequestStatusCompleted, updatedReq.Status)
	})

	t.Run("B端_完成任务_Admin可兜底完成跨站任务", func(t *testing.T) {
		now := time.Now()
		req := &model.ServiceRequest{
			RequestNo:         "REQ-COMPLETE-ADMIN-001",
			UserID:            10,
			ServiceType:       "medical",
			Status:            consts.RequestStatusClaimed,
			ContactName:       "李大爷",
			ContactPhone:      "13900000003",
			Address:           "测试地址-管理员完成",
			StationID:         1,
			CreatedAt:         now,
			UpdatedAt:         now,
			SubmitLocationLat: 39.9,
			SubmitLocationLng: 116.4,
		}
		assert.NoError(t, env.DB.DB.Create(req).Error)

		task := &model.TaskAssignment{
			RequestID: req.ID,
			StationID: 1,
			StaffID:   4,
			Status:    consts.TaskStatusClaimed,
			CreatedAt: now,
			UpdatedAt: now,
			ClaimedAt: now,
		}
		assert.NoError(t, env.DB.DB.Create(task).Error)

		w := testutil.DoRequest(env.Engine, http.MethodPost,
			"/api/v1/b/tasks/"+strconv.FormatInt(task.ID, 10)+"/complete", adminToken, `{"images":["proof-admin.jpg"]}`)
		testutil.AssertOK(t, w)

		var updatedTask model.TaskAssignment
		assert.NoError(t, env.DB.DB.First(&updatedTask, task.ID).Error)
		assert.Equal(t, consts.TaskStatusCompleted, updatedTask.Status)
		assert.NotZero(t, updatedTask.CompletedAt)
	})
}
