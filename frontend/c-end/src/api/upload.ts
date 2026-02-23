import client from './client'

export interface UploadResponse {
  url: string
  key: string
}

export const uploadAPI = {
  /**
   * 上传图片
   * POST /api/v1/c/upload
   * C 端限制：仅 JPEG/PNG/WebP，最大 5MB
   */
  upload: (file: File, module = 'c_end') => {
    const formData = new FormData()
    formData.append('file', file)
    return client.post<any, UploadResponse>(`/c/upload?module=${module}`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }
}
