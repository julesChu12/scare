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

export interface UploadResponse {
    url: string
    filename: string
}

export const commonAPI = {
    /**
     * 文件上传
     * POST /api/v1/c/upload
     */
    upload: (file: File) => {
        const formData = new FormData()
        formData.append('file', file)
        return client.post<any, UploadResponse>('/c/upload', formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        })
    }
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
