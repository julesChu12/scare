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

export const stationAPI = {
    /**
     * 匹配最近的服务站点
     * POST /api/v1/c/stations/match
     * 返回值经 Axios 拦截器解包后直接是 Station 对象
     */
    matchStation: (data: { latitude: number; longitude: number }) => {
        return client.post<any, Station>('/c/stations/match', data)
    }
}
