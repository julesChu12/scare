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

const (
	defaultGeocodeURL        = "https://restapi.amap.com/v3/geocode/geo"
	defaultReverseGeocodeURL = "https://restapi.amap.com/v3/geocode/regeo"
)

// GeocodeService 地址解析服务（高德地图API）
type GeocodeService struct {
	apiKey            string // 高德地图 API Key
	client            *http.Client
	geocodeURL        string
	reverseGeocodeURL string
}

// GeocodeResult 地址解析结果
type GeocodeResult struct {
	Latitude         float64
	Longitude        float64
	FormattedAddress string
}

// NewGeocodeService 创建地理编码服务实例
func NewGeocodeService(apiKey string) *GeocodeService {
	return &GeocodeService{
		apiKey:            apiKey,
		client:            http.DefaultClient,
		geocodeURL:        defaultGeocodeURL,
		reverseGeocodeURL: defaultReverseGeocodeURL,
	}
}

// Geocode 将地址字符串转换为经纬度（地理编码 / Forward Geocoding）
//
// 执行流程：
//   1. 无 API Key（开发环境）→ 返回北京天安门默认坐标，避免接口调用失败
//   2. 构造高德 REST API 请求：GET /v3/geocode/geo?key=...&address=...
//   3. 解析 JSON 响应，提取第一个结果（多结果时取第一个）
//   4. 从 "经度,纬度" 字符串中解析出两个 float64 值
//
// 高德响应格式示例：
//   {
//     "status": "1",
//     "geocodes": [{
//       "location": "116.397470,39.908823",
//       "formatted_address": "北京市东城区天安门广场"
//     }]
//   }
//
// 错误处理：
//   - status != "1"     → 返回 ErrGeocodeAPIError
//   - geocodes 为空     → 返回 ErrGeocodeNotFound（地址无法识别，走人工派单）
//   - geocodes 不为空   → 返回第一条结果（多义性地名取最常用结果）
//
// 在派单链路中的位置：
//   resolveDispatch() 调用 Geocode → 若成功，decision.ServiceLocation 被赋值 → 进入射线法匹配
//   若 ErrGeocodeNotFound → 保持 ServiceLocation=0 → 走人工派单
//
func (s *GeocodeService) Geocode(address string) (*GeocodeResult, error) {
	// 开发环境降级：无 API Key 时返回天安门默认坐标，避免接口调用失败
	if s.apiKey == "" {
		return &GeocodeResult{
			Latitude:         39.908823,
			Longitude:        116.397470,
			FormattedAddress: address,
		}, nil
	}

	// 构造高德地理编码 API 请求
	baseURL := s.getGeocodeURL()
	params := url.Values{}
	params.Add("key", s.apiKey)
	params.Add("address", address)

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := s.getHTTPClient().Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析高德响应 JSON
	var result struct {
		Status   string `json:"status"`
		Info     string `json:"info"`
		Geocodes []struct {
			FormattedAddress string `json:"formatted_address"`
			Location         string `json:"location"` // 格式："经度,纬度"
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

	// 从 "lng,lat" 字符串解析坐标（注意：高德返回的是"经度,纬度"，不是"纬度,经度"）
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

// ReverseGeocode 将经纬度转换为地址信息（逆地理编码 / Reverse Geocoding）
//
// 执行流程：
//   1. 无 API Key（开发环境）→ 返回北京市东城区默认地址
//   2. 构造高德逆地理编码请求：GET /v3/geocode/regeo?key=...&location=lng,lat
//   3. 解析 regeocode 字段，提取省/市/区/标准化地址
//
// 高德响应格式示例：
//   {
//     "status": "1",
//     "regeocode": {
//       "formatted_address": "北京市昌平区霍营街道...",
//       "addressComponent": {
//         "province": "北京市",
//         "city": "北京市",       // 直辖市时与 province 相同
//         "district": "昌平区"
//       }
//     }
//   }
//
// 特殊处理：
//   city 字段在直辖市（北京/上海/天津/重庆）下可能为空数组（高德返回 []），
//   使用 type switch 区分字符串和空数组，统一用 province 填充 city。
//
// 与 Geocode 的关系：
//   Geocode      → 地址 → 坐标（正向，派单时用于获取 ServiceLocation）
//   ReverseGeocode → 坐标 → 地址（逆向，C 端用于显示当前位置）
//
func (s *GeocodeService) ReverseGeocode(lat, lng float64) (*ReverseGeocodeResult, error) {
	// 开发环境降级：无 API Key 时返回默认地址
	if s.apiKey == "" {
		return &ReverseGeocodeResult{
			Province:         "北京市",
			City:             "北京市",
			District:         "东城区",
			FormattedAddress: "北京市东城区天安门广场",
		}, nil
	}

	// 构造高德逆地理编码 API 请求（location 格式为"经度,纬度"）
	baseURL := s.getReverseGeocodeURL()
	params := url.Values{}
	params.Add("key", s.apiKey)
	params.Add("location", fmt.Sprintf("%f,%f", lng, lat))

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := s.getHTTPClient().Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	// 解析高德逆地理编码响应
	var result struct {
		Status    string `json:"status"`
		Info      string `json:"info"`
		Regeocode struct {
			FormattedAddress string `json:"formatted_address"`
			AddressComponent struct {
				Province string `json:"province"`
				City     any    `json:"city"` // 直辖市时可能为空数组
				District string `json:"district"`
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

	// 处理 city 字段（直辖市时 city 为空数组，非直辖市为字符串）
	city := ""
	switch v := result.Regeocode.AddressComponent.City.(type) {
	case string:
		city = v
	default:
		// 直辖市（如北京）city 返回空数组 []，用 province 替代
		city = result.Regeocode.AddressComponent.Province
	}

	return &ReverseGeocodeResult{
		Province:         result.Regeocode.AddressComponent.Province,
		City:             city,
		District:         result.Regeocode.AddressComponent.District,
		FormattedAddress: result.Regeocode.FormattedAddress,
	}, nil
}

func (s *GeocodeService) getHTTPClient() *http.Client {
	if s.client != nil {
		return s.client
	}
	return http.DefaultClient
}

func (s *GeocodeService) getGeocodeURL() string {
	if s.geocodeURL != "" {
		return s.geocodeURL
	}
	return defaultGeocodeURL
}

func (s *GeocodeService) getReverseGeocodeURL() string {
	if s.reverseGeocodeURL != "" {
		return s.reverseGeocodeURL
	}
	return defaultReverseGeocodeURL
}
