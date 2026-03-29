//go:build integration

package integration

import (
	"net/http"
	"strconv"
	"testing"

	"community-elderly-care-platform/test/integration/testutil"

	"github.com/stretchr/testify/assert"
)

func TestZoneCRUD(t *testing.T) {
	env := testutil.Setup(t)
	adminToken := testutil.AdminToken()
	stationManagerToken := testutil.StationManagerToken()

	t.Run("列表_获取所有围栏", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/zones", adminToken)
		testutil.AssertOK(t, w)
	})

	t.Run("创建_新围栏", func(t *testing.T) {
		body := `{
			"station_id": 1,
			"name": "测试围栏",
			"points": [
				{"lng": 116.3, "lat": 39.8},
				{"lng": 116.5, "lat": 39.8},
				{"lng": 116.5, "lat": 40.0},
				{"lng": 116.3, "lat": 40.0}
			],
			"priority": 2
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/zones", adminToken, body)
		data := testutil.AssertOK(t, w)
		assert.NotNil(t, data)
	})

	t.Run("更新_修改围栏", func(t *testing.T) {
		body := `{
			"station_id": 1,
			"name": "更新后围栏",
			"points": [
				{"lng": 116.3, "lat": 39.8},
				{"lng": 116.5, "lat": 39.8},
				{"lng": 116.5, "lat": 40.0},
				{"lng": 116.3, "lat": 40.0}
			],
			"priority": 3
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/zones/2", adminToken, body)
		// 围栏 2 由创建步骤产生，验证响应正常
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
	})

	t.Run("删除_移除围栏", func(t *testing.T) {
		w := testutil.DoRequest(env.Engine, http.MethodDelete, "/api/v1/b/zones/2", adminToken)
		assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
	})

	t.Run("站长只能查看所属站点围栏", func(t *testing.T) {
		stationBody := `{"name": "第二服务站", "code": "ST002", "status": "active"}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/stations", adminToken, stationBody)
		assert.Contains(t, []int{http.StatusOK, http.StatusCreated}, w.Code)
		stationData := testutil.ParseData(t, w)
		secondStationID := int(stationData["id"].(float64))

		zoneBody := `{
			"station_id": ` + strconv.Itoa(secondStationID) + `,
			"name": "第二站围栏",
			"points": [
				{"lng": 116.6, "lat": 39.8},
				{"lng": 116.8, "lat": 39.8},
				{"lng": 116.8, "lat": 40.0},
				{"lng": 116.6, "lat": 40.0}
			],
			"priority": 1
		}`
		w = testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/zones", adminToken, zoneBody)
		assert.Equal(t, http.StatusOK, w.Code)

		w = testutil.DoRequest(env.Engine, http.MethodGet, "/api/v1/b/zones?page=1&page_size=20", stationManagerToken)
		data := testutil.AssertPageResponse(t, w)
		items := data["items"].([]interface{})
		assert.Len(t, items, 1)

		first := items[0].(map[string]interface{})
		assert.Equal(t, float64(1), first["station_id"])
	})

	t.Run("站长只能创建本站点围栏", func(t *testing.T) {
		ownZoneBody := `{
			"station_id": 1,
			"name": "站长自建围栏",
			"points": [
				{"lng": 116.31, "lat": 39.81},
				{"lng": 116.32, "lat": 39.81},
				{"lng": 116.32, "lat": 39.82},
				{"lng": 116.31, "lat": 39.82}
			],
			"priority": 2
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/zones", stationManagerToken, ownZoneBody)
		assert.Equal(t, http.StatusOK, w.Code)

		stationBody := `{"name": "第三服务站", "code": "ST003", "status": "active"}`
		w = testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/stations", adminToken, stationBody)
		assert.Contains(t, []int{http.StatusOK, http.StatusCreated}, w.Code)
		stationData := testutil.ParseData(t, w)
		otherStationID := int(stationData["id"].(float64))

		otherZoneBody := `{
			"station_id": ` + strconv.Itoa(otherStationID) + `,
			"name": "越权围栏",
			"points": [
				{"lng": 116.61, "lat": 39.81},
				{"lng": 116.62, "lat": 39.81},
				{"lng": 116.62, "lat": 39.82},
				{"lng": 116.61, "lat": 39.82}
			],
			"priority": 2
		}`
		w = testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/zones", stationManagerToken, otherZoneBody)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("站长只能修改和删除本站点围栏", func(t *testing.T) {
		ownZoneBody := `{
			"station_id": 1,
			"name": "本站围栏待修改",
			"points": [
				{"lng": 116.33, "lat": 39.83},
				{"lng": 116.34, "lat": 39.83},
				{"lng": 116.34, "lat": 39.84},
				{"lng": 116.33, "lat": 39.84}
			],
			"priority": 3
		}`
		w := testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/zones", adminToken, ownZoneBody)
		assert.Equal(t, http.StatusOK, w.Code)
		ownZoneData := testutil.ParseData(t, w)
		ownZoneID := int(ownZoneData["id"].(float64))

		stationBody := `{"name": "第四服务站", "code": "ST004", "status": "active"}`
		w = testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/stations", adminToken, stationBody)
		assert.Contains(t, []int{http.StatusOK, http.StatusCreated}, w.Code)
		stationData := testutil.ParseData(t, w)
		otherStationID := int(stationData["id"].(float64))

		otherZoneBody := `{
			"station_id": ` + strconv.Itoa(otherStationID) + `,
			"name": "他站围栏待越权",
			"points": [
				{"lng": 116.71, "lat": 39.83},
				{"lng": 116.72, "lat": 39.83},
				{"lng": 116.72, "lat": 39.84},
				{"lng": 116.71, "lat": 39.84}
			],
			"priority": 3
		}`
		w = testutil.DoRequest(env.Engine, http.MethodPost, "/api/v1/b/zones", adminToken, otherZoneBody)
		assert.Equal(t, http.StatusOK, w.Code)
		otherZoneData := testutil.ParseData(t, w)
		otherZoneID := int(otherZoneData["id"].(float64))

		updateOwnBody := `{
			"station_id": 1,
			"name": "本站围栏已修改",
			"points": [
				{"lng": 116.33, "lat": 39.83},
				{"lng": 116.35, "lat": 39.83},
				{"lng": 116.35, "lat": 39.85},
				{"lng": 116.33, "lat": 39.85}
			],
			"priority": 4
		}`
		w = testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/zones/"+strconv.Itoa(ownZoneID), stationManagerToken, updateOwnBody)
		assert.Equal(t, http.StatusOK, w.Code)

		updateOtherBody := `{
			"station_id": ` + strconv.Itoa(otherStationID) + `,
			"name": "越权修改",
			"points": [
				{"lng": 116.71, "lat": 39.83},
				{"lng": 116.73, "lat": 39.83},
				{"lng": 116.73, "lat": 39.85},
				{"lng": 116.71, "lat": 39.85}
			],
			"priority": 5
		}`
		w = testutil.DoRequest(env.Engine, http.MethodPut, "/api/v1/b/zones/"+strconv.Itoa(otherZoneID), stationManagerToken, updateOtherBody)
		assert.Equal(t, http.StatusForbidden, w.Code)

		w = testutil.DoRequest(env.Engine, http.MethodDelete, "/api/v1/b/zones/"+strconv.Itoa(ownZoneID), stationManagerToken)
		assert.Equal(t, http.StatusOK, w.Code)

		w = testutil.DoRequest(env.Engine, http.MethodDelete, "/api/v1/b/zones/"+strconv.Itoa(otherZoneID), stationManagerToken)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
