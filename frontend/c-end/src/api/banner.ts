import client from './client'

export interface Banner {
    id: number
    title: string
    image_url: string
    link_url: string
    sort: number
}

export const bannerAPI = {
    /**
     * 获取轮播图列表
     * GET /api/v1/c/banners
     */
    getBanners: () => {
        return client.get<any, Banner[]>('/c/banners')
    }
}
