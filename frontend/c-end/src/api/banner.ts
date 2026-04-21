import client from './client'

export type LinkType = 'none' | 'news' | 'url'

export interface Banner {
    id: number
    title: string
    image_url: string
    link_type: LinkType
    link_value: string
    sort: number
}

export const bannerAPI = {
    /**
     * 获取轮播图列表
     * GET /api/v1/c/banners?station_id=1
     */
    getBanners: (stationId?: number) => {
        return client.get<any, Banner[]>('/c/banners', { params: { station_id: stationId } })
    }
}
