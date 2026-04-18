package service

import "math"

const earthRadiusMeters = 6371000

// HaversineDistance 使用 Haversine 公式计算两点之间的地球表面距离（单位：米）
//
// 算法原理：
//   Haversine 公式用于计算球面上两点间的最短距离（大圆距离），
//   适用于地球表面经纬度坐标的距离计算。
//
// 公式推导：
//   设地球平均半径 R = 6371000 米
//   两点坐标：(lat1, lng1) 和 (lat2, lng2)，经纬度均以度为单位
//
//   第一步：将经纬度转为弧度
//     lat1Rad = lat1 * π/180
//     lat2Rad = lat2 * π/180
//     Δlat   = (lat2 - lat1) * π/180
//     Δlng   = (lng2 - lng1) * π/180
//
//   第二步：计算半弦长 a
//     a = sin²(Δlat/2) + cos(lat1Rad) * cos(lat2Rad) * sin²(Δlng/2)
//           ^^^^^^^^^^^   ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//           南北分量          东西分量（需乘两点的纬度余弦作投影修正）
//
//   第三步：计算角距离 c
//     c = 2 * atan2(√a, √(1-a))
//              ^^^^^^^^^^^^^^^
//              使用 atan2 代替 arcsin，数值稳定性更好
//
//   第四步：得到球面距离
//     d = R * c
//
// 误差说明：
//   地球并非完美球体（赤道半径略大于极半径），本公式误差约 0.5%，
//   对于社区养老服务场景（公里级距离）完全满足精度要求。
//
func HaversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// 将角度转为弧度（三角函数标准输入）
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180 // 纬度差（弧度）
	deltaLng := (lng2 - lng1) * math.Pi / 180 // 经度差（弧度）

	// 步骤一：计算半弦长 a
	// sin²(Δlat/2)：南北方向的分量贡献
	// cos(lat1)*cos(lat2)*sin²(Δlng/2)：东西方向的分量贡献（纬度越高，投影越小）
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(deltaLng/2)*math.Sin(deltaLng/2)

	// 步骤二：计算角距离 c
	// atan2(√a, √(1-a)) 等价于 2 * asin(√a)，数值稳定性更优
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// 步骤三：乘以地球半径得到实际距离（米）
	return earthRadiusMeters * c
}
