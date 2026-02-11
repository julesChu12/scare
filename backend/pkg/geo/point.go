package geo

// Point 地理坐标点
type Point struct {
	Lat float64 `json:"lat"` // 纬度
	Lng float64 `json:"lng"` // 经度
}

// BoundingBox 外包矩形（用于快速排除）
type BoundingBox struct {
	MinLat float64
	MaxLat float64
	MinLng float64
	MaxLng float64
}

// NewBoundingBox 从点集创建外包矩形
func NewBoundingBox(points []Point) BoundingBox {
	if len(points) == 0 {
		return BoundingBox{}
	}

	minLat, maxLat := points[0].Lat, points[0].Lat
	minLng, maxLng := points[0].Lng, points[0].Lng

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

// Contains 判断点是否在矩形内
func (b BoundingBox) Contains(p Point) bool {
	return p.Lat >= b.MinLat && p.Lat <= b.MaxLat &&
		p.Lng >= b.MinLng && p.Lng <= b.MaxLng
}
