package service

import (
	"encoding/json"
	"sync"

	"community-elderly-care-platform/internal/repository"
	"community-elderly-care-platform/pkg/geo"
)

// GeofenceService 地理围栏服务
//
// 职责：
//   1. 从数据库加载所有活跃围栏，构建内存中的 geo.Engine
//   2. 提供坐标→站点的匹配接口，供派单和站点匹配使用
//   3. 支持热重载（围栏增删改后无需重启服务）
//
// 线程安全：
//   使用 sync.RWMutex 保护 engine 字段。
//   Reload() 持有写锁（独占），Match() 持有读锁（可并发）。
//   读锁期间 engine 指针稳定（Reload 在替换前已完成新 engine 的构建）。
//
// 热重载触发点：
//   - 程序启动时（deps.go 调用一次 Reload）
//   - ZoneService.Create / Update / Delete 后（各自调用 geofenceSvc.Reload）
//
type GeofenceService struct {
	zoneRepo *repository.ZoneRepository
	mu       sync.RWMutex            // 读写锁：写锁保护 engine 替换，读锁保护 Match 读取
	engine   *geo.Engine             // 内存中的围栏匹配引擎，为 nil 时表示未初始化
}

// NewGeofenceService 创建地理围栏服务（需外部调用 Reload 完成初始化）
func NewGeofenceService(zoneRepo *repository.ZoneRepository) *GeofenceService {
	return &GeofenceService{
		zoneRepo: zoneRepo,
	}
}

// Reload 从数据库重新加载所有活跃围栏，构建新的 geo.Engine
//
// 执行流程：
//   1. 从 zoneRepo 查出所有 status=active 的围栏
//   2. 将每个围栏的 JSON Points 字段反序列化为 []geo.Point
//   3. 构建 []geo.Zone（含预计算的 BoundingBox）并排序
//   4. 替换 engine 字段（写锁）
//
// 注意事项：
//   - 反序列化失败（如 JSON 格式错误）的围栏会被静默跳过，不影响其他围栏
//   - Reload 期间 Match 调用仍可正常进行（旧 engine 仍可用）
//   - 新旧 engine 切换是原子的（写锁保护的指针替换）
//
func (s *GeofenceService) Reload() error {
	// 步骤一：从数据库查询所有活跃围栏
	zones, err := s.zoneRepo.ListActive()
	if err != nil {
		return err
	}

	// 步骤二：转换数据格式，构建 geo.Zone 列表
	engineZones := make([]geo.Zone, 0, len(zones))
	for _, zone := range zones {
		var points []geo.Point
		if err := json.Unmarshal([]byte(zone.Points), &points); err != nil {
			continue // JSON 解析失败的围栏跳过，不影响其他围栏
		}
		engineZones = append(engineZones, geo.Zone{
			ID:        zone.ID,
			StationID: zone.StationID,
			Priority:  int(zone.Priority),
			Points:    points,
		})
	}

	// 步骤三：构建新的匹配引擎（按优先级排序、预计算 BBox）
	engine := geo.NewEngine(engineZones)

	// 步骤四：原子替换（写锁保护）
	s.mu.Lock()
	s.engine = engine
	s.mu.Unlock()

	return nil
}

// Match 匹配坐标点所属的服务站点
//
// 线程安全：
//   使用读锁（RLock），允许并发调用。
//   engine 指针在 Reload 时原子替换，读取旧指针不影响正确性。
//
// 内部调用链：
//   Match(lat, lng)
//     → geo.Engine.Match(point)          // 遍历所有围栏
//         → BoundingBox.Contains(point)   // 快速排除
//         → PointInPolygon(point, ...)    // 射线法精确判断
//
func (s *GeofenceService) Match(lat, lng float64) (int64, bool) {
	s.mu.RLock()
	engine := s.engine
	s.mu.RUnlock()

	if engine == nil {
		return 0, false
	}
	return engine.Match(geo.Point{Lat: lat, Lng: lng})
}

