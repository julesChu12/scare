import client from './client'

export interface NewsItem {
  id: number
  title: string
  summary: string
  content: string
  cover_url: string
  type: 'news' | 'notice' | 'activity'
  status: string
  station_id?: number
  author_id?: number
  publish_at: string
  view_count: number
  created_at: string
  updated_at: string
}

export interface NewsListResponse {
  items: NewsItem[]
  total: number
  page: number
  page_size: number
}

export interface NewsListParams {
  page?: number
  page_size?: number
  type?: string
  station_id?: number
}

export const newsAPI = {
  /**
   * 获取新闻列表
   */
  async getList(params: NewsListParams = {}): Promise<NewsListResponse> {
    return client.get('/c/news', { params })
  },

  /**
   * 获取新闻详情
   */
  async getDetail(id: number): Promise<NewsItem> {
    return client.get(`/c/news/${id}`)
  }
}
