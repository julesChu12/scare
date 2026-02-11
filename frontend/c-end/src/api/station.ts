import client from './client'

export interface Station {
    id: number
    name: string
    code?: string
    address: string
    phone: string
    latitude?: number
    longitude?: number
    distance?: number
    service_area?: string
    work_hours?: string
}

export type StationInfo = Station

export interface MatchStationResponse {
    station: Station
}

export const stationAPI = {
    /**
     * 匹配最近的服务站点
     * POST /api/v1/c/stations/match
     */
    matchStation: (data: { latitude: number; longitude: number }) => {
        return client.post<any, MatchStationResponse>('/c/stations/match', data)
    }
}
