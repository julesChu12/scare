import client from './client'

export interface GeocodeResponse {
    latitude: number
    longitude: number
    formatted_address: string
}

export interface ReverseGeocodeResponse {
    province: string
    city: string
    district: string
    address: string
}

export const geocodeAPI = {
    /**
     * 地理编码（地址 -> 坐标）
     * POST /api/v1/c/geocode
     */
    geocode: (data: { address: string }) => {
        return client.post<any, GeocodeResponse>('/c/geocode', data)
    },

    /**
     * 逆地理编码（坐标 -> 地址）
     * GET /api/v1/c/geocode/reverse?lat=xxx&lng=xxx
     */
    reverseGeocode: (params: { latitude: number; longitude: number }) => {
        return client.get<any, ReverseGeocodeResponse>('/c/geocode/reverse', {
            params: {
                lat: params.latitude,
                lng: params.longitude
            }
        })
    }
}
