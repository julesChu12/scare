//go:build integration

package service

import (
	"math"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	amapIntegrationAddress = "北京市东城区天安门广场"
	amapIntegrationLat     = 39.908823
	amapIntegrationLng     = 116.397470
)

// requireAmapIntegrationService 按需启用真实高德联调测试。
func requireAmapIntegrationService(t *testing.T) *GeocodeService {
	t.Helper()

	if testing.Short() {
		t.Skip("short mode: skip amap integration test")
	}
	if os.Getenv("RUN_AMAP_INTEGRATION_TEST") != "1" {
		t.Skip("skip amap integration test: set RUN_AMAP_INTEGRATION_TEST=1")
	}

	key := os.Getenv("AMAP_KEY")
	if key == "" {
		t.Skip("skip amap integration test: AMAP_KEY is required")
	}

	svc := NewGeocodeService(key)
	svc.client = &http.Client{Timeout: 15 * time.Second}
	return svc
}

func assertApproxCoordinate(t *testing.T, got, want, tolerance float64, field string) {
	t.Helper()

	if math.Abs(got-want) > tolerance {
		t.Fatalf("unexpected %s: got %.6f want %.6f +/- %.6f", field, got, want, tolerance)
	}
}

func assertContainsAny(t *testing.T, value string, subs ...string) {
	t.Helper()

	for _, sub := range subs {
		if strings.Contains(value, sub) {
			return
		}
	}
	t.Fatalf("expected %q to contain one of %v", value, subs)
}

func TestGeocodeService_Integration_Geocode_RealAmap(t *testing.T) {
	svc := requireAmapIntegrationService(t)

	result, err := svc.Geocode(amapIntegrationAddress)
	if err != nil {
		t.Fatalf("Geocode returned error: %v", err)
	}
	if result.FormattedAddress == "" {
		t.Fatal("expected formatted address to be non-empty")
	}

	assertApproxCoordinate(t, result.Latitude, amapIntegrationLat, 0.02, "latitude")
	assertApproxCoordinate(t, result.Longitude, amapIntegrationLng, 0.02, "longitude")
	assertContainsAny(t, result.FormattedAddress, "北京", "天安门", "东城")
}

func TestGeocodeService_Integration_ReverseGeocode_RealAmap(t *testing.T) {
	svc := requireAmapIntegrationService(t)

	result, err := svc.ReverseGeocode(amapIntegrationLat, amapIntegrationLng)
	if err != nil {
		t.Fatalf("ReverseGeocode returned error: %v", err)
	}
	if result.Province != "北京市" {
		t.Fatalf("expected province 北京市, got %q", result.Province)
	}
	if result.City == "" {
		t.Fatal("expected city to be non-empty")
	}
	if result.District == "" {
		t.Fatal("expected district to be non-empty")
	}
	if result.FormattedAddress == "" {
		t.Fatal("expected formatted address to be non-empty")
	}

	assertContainsAny(t, result.District, "东城", "西城")
	assertContainsAny(t, result.FormattedAddress, "北京", "天安门", "东城")
}

func TestGeocodeService_Integration_RoundTrip_RealAmap(t *testing.T) {
	svc := requireAmapIntegrationService(t)

	geo, err := svc.Geocode(amapIntegrationAddress)
	if err != nil {
		t.Fatalf("Geocode returned error: %v", err)
	}

	regeo, err := svc.ReverseGeocode(geo.Latitude, geo.Longitude)
	if err != nil {
		t.Fatalf("ReverseGeocode returned error: %v", err)
	}
	if regeo.Province != "北京市" {
		t.Fatalf("expected province 北京市 after round trip, got %q", regeo.Province)
	}
	if regeo.FormattedAddress == "" {
		t.Fatal("expected formatted address after round trip to be non-empty")
	}

	assertContainsAny(t, regeo.FormattedAddress, "北京", "天安门", "东城")
}
