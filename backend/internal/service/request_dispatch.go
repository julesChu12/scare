package service

import (
	"errors"
	"strings"

	"community-elderly-care-platform/internal/repository"
)

// 派单依据常量，标识服务请求的自动派单方式
const (
	DispatchBasisServiceGeofence     = "service_geofence"     // 射线法命中地理围栏，自动派单
	DispatchBasisServiceNearest      = "service_nearest"      // 射线法未命中，Haversine 兜底派最近站点
	DispatchBasisAddressManualReview = "service_address_manual_review" // 无坐标或地址无法解析，需人工审核派单
)

var (
	ErrServiceLocationRequired = errors.New("service location required")
	ErrAddressRequired         = errors.New("address required")
)

// DispatchInput 服务请求派单输入参数
type DispatchInput struct {
	Address          string   // 服务地址文本（用于地理编码获取坐标）
	SubmitLatitude   *float64 // 用户提交坐标 - 纬度（可选）
	SubmitLongitude  *float64 // 用户提交坐标 - 经度（可选）
	ServiceLatitude  *float64 // 服务位置坐标 - 纬度（可选，由地址 geocode 而来）
	ServiceLongitude *float64 // 服务位置坐标 - 经度（可选，由地址 geocode 而来）
	SourceStationID  *int64   // 来源站点ID（可选，跨站点提交场景）
}

// DispatchDecision 派单决策结果
type DispatchDecision struct {
	ResolvedAddress   string  // 规范化后的服务地址（geocode 后取高德返回的标准地址）
	SubmitLatitude    float64 // 原始提交坐标纬度
	SubmitLongitude   float64 // 原始提交坐标经度
	ServiceLatitude   float64 // 解析后的服务位置纬度
	ServiceLongitude  float64 // 解析后的服务位置经度
	SourceStationID   int64   // 有效的来源站点ID（校验后）
	AssignedStationID int64   // 分配的站点ID（0 表示未自动分配）
	DispatchBasis     string  // 派单依据：geofence / nearest / manual_review
	NeedsManualVerify bool    // 是否需要人工审核（地址无法解析时为 true）
}

// resolveDispatch 服务请求派单决策核心函数
//
// 执行流程（优先级从高到低）：
//   1. 解析并校验用户提交的坐标对（SubmitLatitude/Longitude）
//   2. 校验来源站点ID（SourceStationID）
//   3. 若有地址文本，调用 GeocodeService 将地址转换为坐标（ServiceLatitude/Longitude）
//      - Geocode 成功 → 使用解析后的标准地址和坐标
//      - Geocode 失败（网络错误等）→ 降级用 Submit 坐标作为服务位置
//      - Geocode ErrGeocodeNotFound → 无法解析，保持 service_location=0，走人工派单
//   4. 若有 ServiceLocation（坐标有效）→ 进入自动派单：
//      - 射线法围栏匹配 → 命中则派单，basis=geofence
//      - 射线法未命中 → Haversine 兜底，basis=nearest
//   5. 若无 ServiceLocation → 人工派单，basis=manual_review，NeedsManualVerify=true
//
func resolveDispatch(input DispatchInput, stationRepo *repository.StationRepository, geofenceSvc *GeofenceService, geocodeSvc *GeocodeService) (*DispatchDecision, error) {
	address := strings.TrimSpace(input.Address)
	decision := &DispatchDecision{
		ResolvedAddress: address,
	}

	// 步骤一：解析并校验用户提交的坐标对
	if err := assignCoordinatePair(input.SubmitLatitude, input.SubmitLongitude, &decision.SubmitLatitude, &decision.SubmitLongitude); err != nil {
		return nil, err
	}

	// 步骤二：校验来源站点ID有效性
	validSourceStationID, err := resolveSourceStationID(input.SourceStationID, stationRepo)
	if err != nil {
		return nil, err
	}
	decision.SourceStationID = validSourceStationID

	// 步骤三：地址不能为空
	if address == "" {
		return nil, ErrAddressRequired
	}

	// 步骤四：地理编码（地址 → 坐标）
	if geocodeSvc != nil {
		geo, geoErr := geocodeSvc.Geocode(address)
		if geoErr == nil && validCoordinate(geo.Latitude, geo.Longitude) {
			// Geocode 成功：优先使用解析后的标准化地址和坐标
			decision.ServiceLatitude = geo.Latitude
			decision.ServiceLongitude = geo.Longitude
			if geo.FormattedAddress != "" {
				decision.ResolvedAddress = geo.FormattedAddress // 高德标准地址替换用户输入
			}
		} else if geoErr != nil && !errors.Is(geoErr, ErrGeocodeNotFound) {
			// Geocode API 错误（如 Key 平台不匹配、网络异常）→ 降级用 Submit 坐标
			if decision.SubmitLatitude != 0 && decision.SubmitLongitude != 0 {
				decision.ServiceLatitude = decision.SubmitLatitude
				decision.ServiceLongitude = decision.SubmitLongitude
			}
		}
		// ErrGeocodeNotFound：地址无法识别，保持 service_location=0 → 走人工派单
	}

	// 步骤五：有服务坐标 → 自动派单
	hasServiceLocation := decision.ServiceLatitude != 0 && decision.ServiceLongitude != 0
	if hasServiceLocation {
		stationID, matched, matchErr := resolveAssignedStation(decision.ServiceLatitude, decision.ServiceLongitude, stationRepo, geofenceSvc)
		if matchErr != nil {
			return nil, matchErr
		}

		decision.AssignedStationID = stationID
		if matched {
			decision.DispatchBasis = DispatchBasisServiceGeofence
		} else {
			decision.DispatchBasis = DispatchBasisServiceNearest
		}
		return decision, nil
	}

	// 步骤六：无坐标 → 人工派单
	decision.DispatchBasis = DispatchBasisAddressManualReview
	decision.NeedsManualVerify = true
	return decision, nil
}

// assignCoordinatePair 安全解析一对坐标（纬度+经度）
//
// 规则：
//   - 两者都为 nil → 合法的"无坐标"状态，返回 nil
//   - 两者只有一个为 nil → 非法，坐标必须成对出现
//   - 两者都非 nil → 校验坐标有效性，无效返回错误
//
func assignCoordinatePair(lat, lng *float64, outLat, outLng *float64) error {
	if lat == nil && lng == nil {
		return nil
	}
	if lat == nil || lng == nil {
		return ErrInvalidRequest
	}
	if !validCoordinate(*lat, *lng) {
		return ErrInvalidRequest
	}

	*outLat = *lat
	*outLng = *lng
	return nil
}

// resolveSourceStationID 校验并返回有效的来源站点ID
//
// 规则：
//   - sourceStationID 为空或 ≤0 → 返回 0（无来源站点）
//   - 站点不存在或状态非 active → 返回 0（视为无效）
//   - 校验通过 → 返回该站点ID
//
func resolveSourceStationID(sourceStationID *int64, stationRepo *repository.StationRepository) (int64, error) {
	if sourceStationID == nil || *sourceStationID <= 0 || stationRepo == nil {
		return 0, nil
	}

	station, err := stationRepo.GetByID(*sourceStationID)
	if err != nil {
		return 0, nil
	}
	if station.Status != "active" {
		return 0, nil
	}
	return station.ID, nil
}

// resolveAssignedStation 为给定坐标分配服务站点（核心派单逻辑）
//
// 派单策略（两级降级）：
//   第一级：Geofence 射线法匹配
//     调用 geofenceSvc.Match(lat, lng)，内部执行：
//       Engine.Match() → BoundingBox.Contains() + PointInPolygon()
//     若命中（matched=true）→ 直接返回该 StationID
//
//   第二级：Haversine 最近站兜底
//     若射线法未命中任何围栏，从所有活跃站点中遍历计算 Haversine 距离，
//     返回距离最近的站点（matched=false）
//
//   特殊情况：
//     - geofenceSvc 为 nil（未初始化）→ 跳过第一级，直接进入 Haversine
//     - stationRepo 为 nil → 返回错误（无法查找站点）
//     - 无任何活跃站点 → 返回错误
//
// 返回值 (stationID, matched, error)：
//   matched=true  → 射线法命中，basis=geofence
//   matched=false → Haversine 兜底，basis=nearest
//
func resolveAssignedStation(lat, lng float64, stationRepo *repository.StationRepository, geofenceSvc *GeofenceService) (int64, bool, error) {
	// 第一级：射线法围栏匹配
	if geofenceSvc != nil {
		if stationID, matched := geofenceSvc.Match(lat, lng); matched {
			return stationID, true, nil // 命中，返回 StationID
		}
		// 未命中，继续第二级
	}

	// 第二级：Haversine 找最近站点（兜底）
	if stationRepo == nil {
		return 0, false, ErrNoStation
	}

	stations, err := stationRepo.ListActive()
	if err != nil {
		return 0, false, err
	}

	nearest, err := nearestStationByHaversine(stations, lat, lng)
	if err != nil {
		return 0, false, err
	}
	return nearest.ID, false, nil
}
