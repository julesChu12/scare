package geo

import "math"

// PointInPolygon 判断点是否在多边形内（射线法）
func PointInPolygon(point Point, polygon []Point) bool {
	if len(polygon) < 3 {
		return false
	}

	inside := false
	j := len(polygon) - 1

	for i := 0; i < len(polygon); i++ {
		pi := polygon[i]
		pj := polygon[j]

		// 检查是否在边上
		if pointOnSegment(point, pj, pi) {
			return true
		}

		// 射线法判断
		intersects := ((pi.Lat > point.Lat) != (pj.Lat > point.Lat)) &&
			(point.Lng < (pj.Lng-pi.Lng)*(point.Lat-pi.Lat)/(pj.Lat-pi.Lat)+pi.Lng)
		if intersects {
			inside = !inside
		}
		j = i
	}

	return inside
}

// pointOnSegment 判断点是否在线段上
func pointOnSegment(p, a, b Point) bool {
	const epsilon = 1e-9

	// 叉积判断共线
	cross := (b.Lng-a.Lng)*(p.Lat-a.Lat) - (b.Lat-a.Lat)*(p.Lng-a.Lng)
	if math.Abs(cross) > epsilon {
		return false
	}

	// 点积判断在线段内
	dot := (p.Lng-a.Lng)*(p.Lng-b.Lng) + (p.Lat-a.Lat)*(p.Lat-b.Lat)
	return dot <= 0
}
