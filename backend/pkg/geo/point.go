package geo

// Point 地理坐标点
type Point struct {
	Lat float64 `json:"lat"` // 纬度
	Lng float64 `json:"lng"` // 经度
}

// BoundingBox 外包矩形（Axis-Aligned Bounding Box, AABB）
//
// 外包矩形是包含多边形所有顶点的最小 axis-aligned（边与坐标轴平行）矩形。
// 用途：作为 PointInPolygon 的预过滤层。若点不在外包矩形内，则必定不在多边形内，
// 可用 4 次比较（O(1)）快速排除大量候选围栏，避免对每个围栏执行 O(n) 射线法。
//
// 示例：
//   围栏顶点：[(40.0, 116.0), (40.1, 116.1), (40.0, 116.1), (40.1, 116.0)]
//   外包矩形：MinLat=40.0, MaxLat=40.1, MinLng=116.0, MaxLng=116.1
//   点 (40.05, 116.05) → 在矩形内 → 继续射线法精确判断
//   点 (39.9, 116.05)  → 在矩形外 → 直接跳过（必定不在多边形内）
//
type BoundingBox struct {
	MinLat float64 // 最小纬度（矩形南边界）
	MaxLat float64 // 最大纬度（矩形北边界）
	MinLng float64 // 最小经度（矩形西边界）
	MaxLng float64 // 最大经度（矩形东边界）
}

// NewBoundingBox 从多边形顶点集合创建外包矩形
//
// 算法：一次遍历，记录所有顶点的纬度和经度极值。
// 复杂度：O(n)，n 为顶点数量。每个顶点仅比较一次。
//
func NewBoundingBox(points []Point) BoundingBox {
	if len(points) == 0 {
		return BoundingBox{}
	}

	// 以第一个点初始化四个极值
	minLat, maxLat := points[0].Lat, points[0].Lat
	minLng, maxLng := points[0].Lng, points[0].Lng

	// 遍历剩余所有顶点，实时更新极值
	for _, p := range points[1:] {
		if p.Lat < minLat {
			minLat = p.Lat
		}
		if p.Lat > maxLat {
			maxLat = p.Lat
		}
		if p.Lng < minLng {
			minLng = p.Lng
		}
		if p.Lng > maxLng {
			maxLng = p.Lng
		}
	}

	return BoundingBox{
		MinLat: minLat,
		MaxLat: maxLat,
		MinLng: minLng,
		MaxLng: maxLng,
	}
}

// Contains 判断点是否在矩形内（闭区间，含边界）
//
// 判断条件：MinLat ≤ point.Lat ≤ MaxLat 且 MinLng ≤ point.Lng ≤ MaxLng
// 注意：包含边界，与 PointInPolygon 的边命中行为保持一致。
// 效率：4 次比较，O(1)。当围栏数量多且分布稀疏时，BBox 过滤可跳过大部分围栏。
//
func (b BoundingBox) Contains(p Point) bool {
	return p.Lat >= b.MinLat && p.Lat <= b.MaxLat &&
		p.Lng >= b.MinLng && p.Lng <= b.MaxLng
}
