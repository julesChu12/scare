package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	ErrGeocodeNotFound        = errors.New("address not found")
	ErrGeocodeAPIError        = errors.New("geocode api error")
	ErrReverseGeocodeNotFound = errors.New("location not found")
)

// GeocodeService 地址解析服务（高德地图API）
type GeocodeService struct {
	apiKey string // 高德地图 API Key
}

// GeocodeResult 地址解析结果
type GeocodeResult struct {
	Latitude         float64
	Longitude        float64
	FormattedAddress string
}

func NewGeocodeService(apiKey string) *GeocodeService {
	return &GeocodeService{apiKey: apiKey}
}

// Geocode 将地址字符串转换为经纬度
func (s *GeocodeService) Geocode(address string) (*GeocodeResult, error) {
	// 如果没有配置 API Key，返回默认值
	if s.apiKey == "" {
		// 开发环境：返回默认经纬度（北京天安门）
		return &GeocodeResult{
			Latitude:         39.908823,
			Longitude:        116.397470,
			FormattedAddress: address,
		}, nil
	}

	// 调用高德地图地理编码 API
	baseURL := "https://restapi.amap.com/v3/geocode/geo"
	params := url.Values{}
	params.Add("key", s.apiKey)
	params.Add("address", address)

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var result struct {
		Status   string `json:"status"`
		Info     string `json:"info"`
		Geocodes []struct {
			FormattedAddress string `json:"formatted_address"`
			Location         string `json:"location"` // "经度,纬度"
		} `json:"geocodes"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("%w: %s", ErrGeocodeAPIError, result.Info)
	}

	if len(result.Geocodes) == 0 {
		return nil, ErrGeocodeNotFound
	}

	// 解析经纬度（格式：经度,纬度）
	var lng, lat float64
	_, err = fmt.Sscanf(result.Geocodes[0].Location, "%f,%f", &lng, &lat)
	if err != nil {
		return nil, err
	}

	return &GeocodeResult{
		Latitude:         lat,
		Longitude:        lng,
		FormattedAddress: result.Geocodes[0].FormattedAddress,
	}, nil
}

// ReverseGeocodeResult 逆地理编码结果
type ReverseGeocodeResult struct {
	Province         string
	City             string
	District         string
	FormattedAddress string
}

// ReverseGeocode 将经纬度转换为地址信息
func (s *GeocodeService) ReverseGeocode(lat, lng float64) (*ReverseGeocodeResult, error) {
	// 如果没有配置 API Key，返回默认值
	if s.apiKey == "" {
		// 开发环境：返回默认地址
		return &ReverseGeocodeResult{
			Province:         "北京市",
			City:             "北京市",
			District:         "东城区",
			FormattedAddress: "北京市东城区天安门广场",
		}, nil
	}

	// 调用高德地图逆地理编码 API
	baseURL := "https://restapi.amap.com/v3/geocode/regeo"
	params := url.Values{}
	params.Add("key", s.apiKey)
	params.Add("location", fmt.Sprintf("%f,%f", lng, lat)) // 高德格式：经度,纬度

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	// 调试日志
	fmt.Printf("[Geocode] lat=%f lng=%f API Response: %s\n", lat, lng, string(body))

	// 解析响应
	var result struct {
		Status    string `json:"status"`
		Info      string `json:"info"`
		Regeocode struct {
			FormattedAddress string `json:"formatted_address"`
			AddressComponent struct {
				Province string      `json:"province"`
				City     interface{} `json:"city"`     // 可能是字符串或空数组
				District string      `json:"district"`
			} `json:"addressComponent"`
		} `json:"regeocode"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("%w: %s", ErrGeocodeAPIError, result.Info)
	}

	if result.Regeocode.FormattedAddress == "" {
		return nil, ErrReverseGeocodeNotFound
	}

	// 处理 city 字段（直辖市时可能为空数组）
	city := ""
	switch v := result.Regeocode.AddressComponent.City.(type) {
	case string:
		city = v
	default:
		// 直辖市时 city 为空数组，使用 province
		city = result.Regeocode.AddressComponent.Province
	}

	return &ReverseGeocodeResult{
		Province:         result.Regeocode.AddressComponent.Province,
		City:             city,
		District:         result.Regeocode.AddressComponent.District,
		FormattedAddress: result.Regeocode.FormattedAddress,
	}, nil
}
