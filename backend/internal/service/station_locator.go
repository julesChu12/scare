package service

import "community-elderly-care-platform/internal/dao/model"

// nearestStationByHaversine 根据 Haversine 球面距离查找最近的站点
//
// 算法：线性遍历所有站点，逐个计算 Haversine 距离，维护最小值。
// 复杂度：O(n)，n = 站点数量。
//
// 适用场景（射线法的兜底方案）：
//   - 用户坐标落在所有围栏外部，射线法无命中
//   - 系统中围栏数据尚未配置
//   - 需要从多个候选站点中选出最近的一个
//
// 与射线法的关系：
//   射线法命中 → 直接返回命中的 StationID（精确匹配，优先级最高）
//   射线法未命中 → Haversine 兜底找最近站点（模糊匹配，作为保底策略）
//
// 特殊情况：
//   - stations 为空或全为 nil → 返回 ErrNoStation
//   - 所有站点距离相同（极端情况）→ 返回遍历中最后一个（无确定性保证）
//
func nearestStationByHaversine(stations []*model.ServiceStation, lat, lng float64) (*model.ServiceStation, error) {
	var nearest *model.ServiceStation
	minDistance := 0.0

	// 线性扫描：计算每个站点到目标点的 Haversine 距离，维护最近站点
	for _, station := range stations {
		if station == nil {
			continue
		}

		// HaversineDistance 计算地球表面两点间的球面距离（单位：米）
		distance := HaversineDistance(lat, lng, station.Latitude, station.Longitude)

		// 首次迭代：nearest == nil，直接记录第一个有效站点
		// 后续迭代：distance < minDistance 时更新最近站点
		if nearest == nil || distance < minDistance {
			nearest = station
			minDistance = distance
		}
	}

	if nearest == nil {
		return nil, ErrNoStation
	}

	return nearest, nil
}
