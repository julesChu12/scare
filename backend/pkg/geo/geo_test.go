package geo

import "testing"

func TestPointInPolygon(t *testing.T) {
	// 正方形：(0,0) -> (0,1) -> (1,1) -> (1,0)
	polygon := []Point{
		{Lat: 0, Lng: 0},
		{Lat: 0, Lng: 1},
		{Lat: 1, Lng: 1},
		{Lat: 1, Lng: 0},
	}

	tests := []struct {
		name   string
		point  Point
		expect bool
	}{
		{"内部点", Point{Lat: 0.5, Lng: 0.5}, true},
		{"外部点", Point{Lat: 1.5, Lng: 0.5}, false},
		{"边界点", Point{Lat: 1.0, Lng: 0.5}, true},
		{"顶点", Point{Lat: 0, Lng: 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PointInPolygon(tt.point, polygon)
			if got != tt.expect {
				t.Errorf("PointInPolygon(%v) = %v, want %v", tt.point, got, tt.expect)
			}
		})
	}
}

func TestEngineMatchPriority(t *testing.T) {
	// 两个重叠的区域，测试优先级匹配
	zones := []Zone{
		{
			ID:        1,
			StationID: 100,
			Priority:  5,
			Points: []Point{
				{Lat: 0, Lng: 0},
				{Lat: 0, Lng: 2},
				{Lat: 2, Lng: 2},
				{Lat: 2, Lng: 0},
			},
		},
		{
			ID:        2,
			StationID: 200,
			Priority:  10, // 更高优先级
			Points: []Point{
				{Lat: 1, Lng: 1},
				{Lat: 1, Lng: 3},
				{Lat: 3, Lng: 3},
				{Lat: 3, Lng: 1},
			},
		},
	}

	engine := NewEngine(zones)

	// 重叠区域内的点应匹配高优先级
	stationID, matched := engine.Match(Point{Lat: 1.5, Lng: 1.5})
	if !matched {
		t.Fatal("expected match")
	}
	if stationID != 200 {
		t.Errorf("expected station 200, got %d", stationID)
	}

	// 只在低优先级区域内的点
	stationID, matched = engine.Match(Point{Lat: 0.5, Lng: 0.5})
	if !matched {
		t.Fatal("expected match")
	}
	if stationID != 100 {
		t.Errorf("expected station 100, got %d", stationID)
	}

	// 两个区域外的点
	_, matched = engine.Match(Point{Lat: 5, Lng: 5})
	if matched {
		t.Error("expected no match")
	}
}

func TestBoundingBox(t *testing.T) {
	points := []Point{
		{Lat: 1, Lng: 2},
		{Lat: 3, Lng: 4},
		{Lat: 2, Lng: 1},
	}

	box := NewBoundingBox(points)

	if box.MinLat != 1 || box.MaxLat != 3 {
		t.Errorf("Lat range: got [%f, %f], want [1, 3]", box.MinLat, box.MaxLat)
	}
	if box.MinLng != 1 || box.MaxLng != 4 {
		t.Errorf("Lng range: got [%f, %f], want [1, 4]", box.MinLng, box.MaxLng)
	}

	// 测试 Contains
	if !box.Contains(Point{Lat: 2, Lng: 2}) {
		t.Error("expected point inside box")
	}
	if box.Contains(Point{Lat: 0, Lng: 2}) {
		t.Error("expected point outside box")
	}
}
