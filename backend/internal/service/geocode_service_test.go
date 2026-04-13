package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeocodeService_Geocode_ReturnsDefaultWhenAPIKeyMissing(t *testing.T) {
	svc := NewGeocodeService("")

	result, err := svc.Geocode("北京市东城区")
	if err != nil {
		t.Fatalf("Geocode returned error: %v", err)
	}
	if result.Latitude != 39.908823 || result.Longitude != 116.397470 {
		t.Fatalf("unexpected default coordinates: %+v", result)
	}
	if result.FormattedAddress != "北京市东城区" {
		t.Fatalf("expected formatted address to echo input, got %q", result.FormattedAddress)
	}
}

func TestGeocodeService_Geocode_ParsesAmapResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/geocode/geo" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Fatalf("unexpected key: %q", got)
		}
		if got := r.URL.Query().Get("address"); got != "北京市东城区" {
			t.Fatalf("unexpected address: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"1","info":"OK","geocodes":[{"formatted_address":"北京市东城区天安门广场","location":"116.397470,39.908823"}]}`))
	}))
	defer server.Close()

	svc := NewGeocodeService("test-key")
	svc.client = server.Client()
	svc.geocodeURL = server.URL + "/v3/geocode/geo"

	result, err := svc.Geocode("北京市东城区")
	if err != nil {
		t.Fatalf("Geocode returned error: %v", err)
	}
	if result.Latitude != 39.908823 || result.Longitude != 116.397470 {
		t.Fatalf("unexpected coordinates: %+v", result)
	}
	if result.FormattedAddress != "北京市东城区天安门广场" {
		t.Fatalf("unexpected formatted address: %q", result.FormattedAddress)
	}
}

func TestGeocodeService_Geocode_ReturnsNotFoundWhenAmapHasNoMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"1","info":"OK","geocodes":[]}`))
	}))
	defer server.Close()

	svc := NewGeocodeService("test-key")
	svc.client = server.Client()
	svc.geocodeURL = server.URL + "/v3/geocode/geo"

	_, err := svc.Geocode("不存在的地址")
	if !errors.Is(err, ErrGeocodeNotFound) {
		t.Fatalf("expected ErrGeocodeNotFound, got %v", err)
	}
}

func TestGeocodeService_ReverseGeocode_ReturnsDefaultWhenAPIKeyMissing(t *testing.T) {
	svc := NewGeocodeService("")

	result, err := svc.ReverseGeocode(39.908823, 116.397470)
	if err != nil {
		t.Fatalf("ReverseGeocode returned error: %v", err)
	}
	if result.Province != "北京市" || result.City != "北京市" || result.District != "东城区" {
		t.Fatalf("unexpected default reverse geocode result: %+v", result)
	}
}

func TestGeocodeService_ReverseGeocode_ParsesAmapResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/geocode/regeo" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Fatalf("unexpected key: %q", got)
		}
		if got := r.URL.Query().Get("location"); got != "116.397470,39.908823" {
			t.Fatalf("unexpected location: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"1","info":"OK","regeocode":{"formatted_address":"北京市东城区天安门广场","addressComponent":{"province":"北京市","city":[],"district":"东城区"}}}`))
	}))
	defer server.Close()

	svc := NewGeocodeService("test-key")
	svc.client = server.Client()
	svc.reverseGeocodeURL = server.URL + "/v3/geocode/regeo"

	result, err := svc.ReverseGeocode(39.908823, 116.397470)
	if err != nil {
		t.Fatalf("ReverseGeocode returned error: %v", err)
	}
	if result.Province != "北京市" || result.City != "北京市" || result.District != "东城区" {
		t.Fatalf("unexpected reverse geocode result: %+v", result)
	}
	if result.FormattedAddress != "北京市东城区天安门广场" {
		t.Fatalf("unexpected formatted address: %q", result.FormattedAddress)
	}
}

func TestGeocodeService_ReverseGeocode_ReturnsNotFoundWhenAmapHasNoAddress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"1","info":"OK","regeocode":{"formatted_address":"","addressComponent":{"province":"北京市","city":[],"district":"东城区"}}}`))
	}))
	defer server.Close()

	svc := NewGeocodeService("test-key")
	svc.client = server.Client()
	svc.reverseGeocodeURL = server.URL + "/v3/geocode/regeo"

	_, err := svc.ReverseGeocode(39.908823, 116.397470)
	if !errors.Is(err, ErrReverseGeocodeNotFound) {
		t.Fatalf("expected ErrReverseGeocodeNotFound, got %v", err)
	}
}
