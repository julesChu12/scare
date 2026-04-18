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
//
// 初始化流程：
//   1. 过滤掉顶点数 < 3 的非法围栏（无法构成多边形）
//   2. 为每个有效围栏预计算外包矩形（BoundingBox），用于快速排除
//   3. 按优先级降序排列，匹配时优先命中高优先级围栏
//
// 注意：
//   围栏优先级相同时，Engine 按数组顺序遍历（无确定性保证），
//   因此需在业务层保证同区域围栏优先级互不相同。
//
func NewEngine(zones []Zone) *Engine {
	filtered := make([]Zone, 0, len(zones))

	for _, zone := range zones {
		if len(zone.Points) < 3 {
			continue // 至少需要3个点构成多边形
		}
		zone.Box = NewBoundingBox(zone.Points) // 预计算外包矩形，加速 Contains 判断
		filtered = append(filtered, zone)
	}

	// 按优先级降序排列，高优先级围栏优先匹配
	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].Priority > filtered[j].Priority
	})

	return &Engine{zones: filtered}
}

// Match 匹配坐标点所属的服务站点
// 返回匹配到的 StationID 和是否匹配成功
//
// 匹配策略：两级判断
//   第一级（BoundingBox 快速排除）：
//     用外包矩形做初筛。外包矩形为围栏所有顶点的 Lat/Lng 极值构成的 axis-aligned 矩形。
//     若点不在矩形内，则必定不在多边形内，可直接跳过。
//     这步复杂度为 O(1)，大量围栏时能快速过滤掉绝大多数候选。
//
//   第二级（PointInPolygon 精确判断）：
//     点在外包矩形内，再用射线法做精确判断。
//     射线法复杂度 O(n)，n 为围栏顶点数。
//
// 整体复杂度：
//   最坏 O(m·n)，m=围栏数，n=平均顶点数；BBox 排除使平均复杂度远优于最坏情况。
//   优先级相同时，返回首个命中的 StationID。
//
func (e *Engine) Match(point Point) (int64, bool) {
	for _, zone := range e.zones {
		// 第一级：外包矩形快速排除，不命中则直接跳过
		if !zone.Box.Contains(point) {
			continue
		}
		// 第二级：射线法精确判断点在多边形内
		if PointInPolygon(point, zone.Points) {
			return zone.StationID, true
		}
	}
	return 0, false
}

