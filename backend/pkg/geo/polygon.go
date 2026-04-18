package geo

import "math"

// PointInPolygon 判断点是否在多边形内（射线法 / Ray Casting Algorithm）
//
// 算法原理：
//   从待测点向右（正Lng方向）发射一条水平射线，统计射线与多边形所有边的交点个数。
//
//   - 交点为奇数  →  点在多边形内部（inside = true）
//   - 交点为偶数  →  点在多边形外部（inside = false）
//
// 特殊情况处理：
//   - 点落在多边形边上（包括顶点）→ 视为命中，返回 true
//   - 射线与边顶点重合 → 仅统计一次（通过 j=i 的交错处理避免重复）
//
// 示意（交点数 = 3，为奇数，点在内部）：
//      \          /
//       \   P   /        P 为待测点
//        \  |  /         水平向右的虚线为射线
//         \ | /          与 3 条边相交（奇数），判定为内部
//          \|/
// ---------+----------→ Lng 方向（射线）
//
func PointInPolygon(point Point, polygon []Point) bool {
	if len(polygon) < 3 {
		return false
	}

	inside := false
	// j 初始化为最后一个顶点的索引，与 i 形成"前一条边"
	j := len(polygon) - 1

	for i := 0; i < len(polygon); i++ {
		pi := polygon[i] // 当前边的起点
		pj := polygon[j] // 当前边的终点（上一轮迭代的顶点）

		// 第一步：检查点是否恰好落在当前边线段上
		if pointOnSegment(point, pj, pi) {
			return true
		}

		// 第二步：判断射线是否与当前边相交
		//
		// 条件一：(pi.Lat > point.Lat) != (pj.Lat > point.Lat)
		//   确保射线（水平线 y=point.Lat）能穿过这条边。
		//   即边的两个端点必须在射线的两侧（一个高于点，一个低于点）。
		//   注意：故意用 ">" 而不是 ">="，排除射线经过顶点的情况，
		//   防止同一条边被统计两次（重复统计会导致奇偶判断出错）。
		//
		// 条件二：point.Lng < (pj.Lng-pi.Lng)*(point.Lat-pi.Lat)/(pj.Lat-pi.Lat)+pi.Lng
		//   确保交点位于射线的正方向（点右侧），
		//   即交点的经度必须大于待测点的经度。
		//   公式由边两点式直线方程推导而来：
		//     (Lng - pi.Lng) / (pj.Lng - pi.Lng) = (Lat - pi.Lat) / (pj.Lat - pi.Lat)
		//   解出 Lat=point.Lat 时的 Lng 值。
		//
		intersects := ((pi.Lat > point.Lat) != (pj.Lat > point.Lat)) &&
			(point.Lng < (pj.Lng-pi.Lng)*(point.Lat-pi.Lat)/(pj.Lat-pi.Lat)+pi.Lng)
		if intersects {
			// 每穿过一条边，奇偶翻转一次
			inside = !inside
		}

		// 当前顶点成为下一条边的终点
		j = i
	}

	return inside
}

// pointOnSegment 判断点 P 是否在线段 AB 上（包含端点）
//
// 两步判断法：
//   1. 叉积（外积）为零 → 向量 AP 与 AB 共线（点在直线延长线上）
//   2. 点积 ≤ 0        → 点投影落在 AB 线段区间内（含端点）
//
// 叉积几何含义：
//   cross = (B-A) × (P-A)
//   cross > 0：P 在 AB 左侧（逆时针方向）
//   cross < 0：P 在 AB 右侧（顺时针方向）
//   cross = 0：三点共线
//
// 点积几何含义：
//   dot = (P-A) · (P-B)
//   dot > 0：P 在 A、B 两端点的外侧（投影在线段延长线上）
//   dot = 0：P 与 A 或 B 重合
//   dot < 0：P 的投影落在 A、B 之间
//
// 两者同时满足 → P 在线段 AB 上
//
func pointOnSegment(p, a, b Point) bool {
	const epsilon = 1e-9

	// 步骤一：共线检测（叉积判断）
	// (b.Lng-a.Lng)*(p.Lat-a.Lat) - (b.Lat-a.Lat)*(p.Lng-a.Lng)
	cross := (b.Lng-a.Lng)*(p.Lat-a.Lat) - (b.Lat-a.Lat)*(p.Lng-a.Lng)
	if math.Abs(cross) > epsilon {
		return false
	}

	// 步骤二：投影区间检测（点积判断）
	// (p.Lng-a.Lng)*(p.Lng-b.Lng) + (p.Lat-a.Lat)*(p.Lat-b.Lat)
	dot := (p.Lng-a.Lng)*(p.Lng-b.Lng) + (p.Lat-a.Lat)*(p.Lat-b.Lat)
	// dot <= 0 包含端点重合的情况
	return dot <= 0
}
