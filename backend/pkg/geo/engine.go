package geo

import "sort"

// Zone 地理围栏区域
type Zone struct {
	ID        int64
	StationID int64
	Priority  int
	Points    []Point
	Box       BoundingBox
}

// Engine 围栏匹配引擎
type Engine struct {
	zones []Zone
}

// NewEngine 创建围栏匹配引擎
func NewEngine(zones []Zone) *Engine {
	filtered := make([]Zone, 0, len(zones))

	for _, zone := range zones {
		if len(zone.Points) < 3 {
			continue // 至少需要3个点构成多边形
		}
		zone.Box = NewBoundingBox(zone.Points)
		filtered = append(filtered, zone)
	}

	// 按优先级降序排列
	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].Priority > filtered[j].Priority
	})

	return &Engine{zones: filtered}
}

// Match 匹配坐标点所属的服务站点
// 返回匹配到的 StationID 和是否匹配成功
func (e *Engine) Match(point Point) (int64, bool) {
	for _, zone := range e.zones {
		// 先用外包矩形快速排除
		if !zone.Box.Contains(point) {
			continue
		}
		// 精确判断点是否在多边形内
		if PointInPolygon(point, zone.Points) {
			return zone.StationID, true
		}
	}
	return 0, false
}

// Zones 获取所有围栏区域
func (e *Engine) Zones() []Zone {
	return e.zones
}
