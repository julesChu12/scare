package service

import "community-elderly-care-platform/internal/dao/model"

func nearestStationByHaversine(stations []*model.ServiceStation, lat, lng float64) (*model.ServiceStation, error) {
	var nearest *model.ServiceStation
	minDistance := 0.0

	for _, station := range stations {
		if station == nil {
			continue
		}

		distance := HaversineDistance(lat, lng, station.Latitude, station.Longitude)
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
