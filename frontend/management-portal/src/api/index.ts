import request from '@/utils/request'
import type {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  RefreshResponse,
  MeResponse,
  Task,
  PaginationParams,
  TaskListParams,
  PaginationData,
  ServiceRequest,
  CreateServiceRequest,
  UpdateServiceRequest,
  CompleteTaskRequest,
  Station,
  StationRequest,
  Zone,
  ZoneRequest,
  User,
  UserCreateRequest,
  UserUpdateRequest,
  UpdateUserRolesRequest,
  UpdateRolePermissionsRequest,
  PermissionNode,
  Notification,
  Menu,
  MenuRequest,
  BatchUpdateSortRequest,
  Banner,
  BannerRequest,
  DashboardStats,
  TaskStatsData,
  RequestStatsData,
  TodayStatsData,
  News,
  NewsRequest,
  NewsListParams,
  UpdateUserIdentitiesRequest,
  OverviewStatsData,
  ServiceTypeStatsData,
  TrendItemData,
  EfficiencyStatsData,
  StaffRankingItemData,
  GenerateReportRequest,
  ReportData,
} from '@/types/api'

/**
 * B端认证相关 API
 */
export const authApi = {
  /**
   * 登录
   */
  login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return request.post('/b/auth/login', data).then((res) => res.data)
  },

  /**
   * 刷新 Token
   */
  refreshToken(refreshToken: string): Promise<ApiResponse<RefreshResponse>> {
    return request.post('/b/auth/refresh', { refresh_token: refreshToken }).then((res) => res.data)
  },

  /**
   * 获取当前用户信息（含权限）
   */
  getCurrentUser(): Promise<ApiResponse<MeResponse>> {
    return request.get('/b/auth/me').then((res) => res.data)
  },

  /**
   * 登出
   */
  logout(): Promise<ApiResponse<null>> {
    return request.post('/b/auth/logout').then((res) => res.data)
  },
}

/**
 * 任务相关 API
 */
export const taskApi = {
  /**
   * 获取所有任务列表（管理端通用，复用 pool 接口）
   */
  getTasks(params?: TaskListParams): Promise<ApiResponse<PaginationData<Task>>> {
    return request.get('/b/tasks/pool', { params }).then((res) => res.data)
  },

  /**
   * 获取任务池列表（待认领任务）
   */
  getTaskPool(params?: PaginationParams): Promise<ApiResponse<PaginationData<Task>>> {
    return request.get('/b/tasks/pool', { params: { ...params, status: 'dispatched' } }).then((res) => res.data)
  },

  /**
   * 获取我的任务列表
   */
  getMyTasks(params?: TaskListParams): Promise<ApiResponse<PaginationData<Task>>> {
    return request.get('/b/tasks/my', { params }).then((res) => res.data)
  },

  /**
   * 获取任务详情
   */
  getTask(taskId: number): Promise<ApiResponse<Task>> {
    return request.get(`/b/tasks/${taskId}`).then((res) => res.data)
  },

  /**
   * 认领任务
   */
  claimTask(taskId: number): Promise<ApiResponse<Task>> {
    return request.post(`/b/tasks/${taskId}/claim`).then((res) => res.data)
  },

  /**
   * 完成任务
   */
  completeTask(taskId: number, data: CompleteTaskRequest): Promise<ApiResponse<Task>> {
    return request.post(`/b/tasks/${taskId}/complete`, data).then((res) => res.data)
  },

  /**
   * 转派/指派任务
   */
  transferTask(taskId: number, staffId: number): Promise<ApiResponse<Task>> {
    return request.post(`/b/tasks/${taskId}/transfer`, { staff_id: staffId }).then((res) => res.data)
  },
}

/**
 * 服务需求相关 API
 */
export const requestApi = {
  /**
   * 获取服务需求详情
   */
  getRequest(requestId: number): Promise<ApiResponse<ServiceRequest>> {
    return request.get(`/b/requests/${requestId}`).then((res) => res.data)
  },

  /**
   * 获取服务需求列表
   */
  getRequests(params?: PaginationParams): Promise<ApiResponse<PaginationData<ServiceRequest>>> {
    return request.get('/b/requests', { params }).then((res) => res.data)
  },

  /**
   * 创建服务需求
   */
  createRequest(data: CreateServiceRequest): Promise<ApiResponse<ServiceRequest>> {
    return request.post('/b/requests', data).then((res) => res.data)
  },

  /**
   * 更新服务需求
   */
  updateRequest(requestId: number, data: UpdateServiceRequest): Promise<ApiResponse<ServiceRequest>> {
    return request.put(`/b/requests/${requestId}`, data).then((res) => res.data)
  },

  /**
   * 删除服务需求
   */
  deleteRequest(requestId: number): Promise<ApiResponse<null>> {
    return request.delete(`/b/requests/${requestId}`).then((res) => res.data)
  },
}

/**
 * 服务站点相关 API
 */
export const stationApi = {
  /**
   * 获取站点列表
   */
  getStations(params?: PaginationParams): Promise<ApiResponse<PaginationData<Station>>> {
    return request.get('/b/stations', { params }).then((res) => res.data)
  },

  /**
   * 获取站点详情
   */
  getStation(stationId: number): Promise<ApiResponse<Station>> {
    return request.get(`/b/stations/${stationId}`).then((res) => res.data)
  },

  /**
   * 创建站点
   */
  createStation(data: StationRequest): Promise<ApiResponse<Station>> {
    return request.post('/b/stations', data).then((res) => res.data)
  },

  /**
   * 更新站点
   */
  updateStation(stationId: number, data: StationRequest): Promise<ApiResponse<Station>> {
    return request.put(`/b/stations/${stationId}`, data).then((res) => res.data)
  },

  /**
   * 删除站点
   */
  deleteStation(stationId: number): Promise<ApiResponse<null>> {
    return request.delete(`/b/stations/${stationId}`).then((res) => res.data)
  },
}

/**
 * 服务围栏相关 API
 */
export const zoneApi = {
  /**
   * 获取围栏列表
   */
  getZones(params?: PaginationParams): Promise<ApiResponse<PaginationData<Zone>>> {
    return request.get('/b/zones', { params }).then((res) => res.data)
  },

  /**
   * 创建围栏
   */
  createZone(data: ZoneRequest): Promise<ApiResponse<Zone>> {
    return request.post('/b/zones', data).then((res) => res.data)
  },

  /**
   * 更新围栏
   */
  updateZone(zoneId: number, data: ZoneRequest): Promise<ApiResponse<Zone>> {
    return request.put(`/b/zones/${zoneId}`, data).then((res) => res.data)
  },

  /**
   * 删除围栏
   */
  deleteZone(zoneId: number): Promise<ApiResponse<null>> {
    return request.delete(`/b/zones/${zoneId}`).then((res) => res.data)
  },
}

/**
 * 用户管理相关 API
 */
export const userApi = {
  /**
   * 获取用户列表
   */
  getUsers(params?: PaginationParams): Promise<ApiResponse<PaginationData<User>>> {
    return request.get('/b/users', { params }).then((res) => res.data)
  },

  /**
   * 获取用户详情
   */
  getUser(userId: number): Promise<ApiResponse<User>> {
    return request.get(`/b/users/${userId}`).then((res) => res.data)
  },

  /**
   * 创建用户
   */
  createUser(data: UserCreateRequest): Promise<ApiResponse<User>> {
    return request.post('/b/users', data).then((res) => res.data)
  },

  /**
   * 更新用户
   */
  updateUser(userId: number, data: UserUpdateRequest): Promise<ApiResponse<User>> {
    return request.put(`/b/users/${userId}`, data).then((res) => res.data)
  },

  /**
   * 更新用户角色
   */
  updateUserRoles(userId: number, data: UpdateUserRolesRequest): Promise<ApiResponse<{ user_id: number; roles: string[]; tokens_revoked: boolean }>> {
    return request.put(`/b/users/${userId}/roles`, data).then((res) => res.data)
  },

  /**
   * 更新用户身份
   */
  updateUserIdentities(userId: number, data: UpdateUserIdentitiesRequest): Promise<ApiResponse<{ user_id: number; identities: string[] }>> {
    return request.put(`/b/users/${userId}/identities`, data).then((res) => res.data)
  },
}

/**
 * 权限管理相关 API
 */
export const permissionApi = {
  /**
   * 获取权限树
   */
  getPermissionTree(): Promise<ApiResponse<{ tree: PermissionNode[] }>> {
    return request.get('/b/permissions/tree').then((res) => res.data)
  },

  /**
   * 获取角色权限
   */
  getRolePermissions(role: string): Promise<ApiResponse<{ role: string; permissions: string[] }>> {
    return request.get(`/b/roles/${role}/permissions`).then((res) => res.data)
  },

  /**
   * 更新角色权限
   */
  updateRolePermissions(role: string, data: UpdateRolePermissionsRequest): Promise<ApiResponse<{ role: string; affected_users: number; tokens_revoked: boolean }>> {
    return request.put(`/b/roles/${role}/permissions`, data).then((res) => res.data)
  },
}

/**
 * 通知相关 API
 */
export const notificationApi = {
  /**
   * 获取通知列表
   */
  getNotifications(params?: PaginationParams): Promise<ApiResponse<PaginationData<Notification>>> {
    return request.get('/b/notifications', { params }).then((res) => res.data)
  },

  /**
   * 标记通知为已读
   */
  markAsRead(notificationId: number): Promise<ApiResponse<null>> {
    return request.post(`/b/notifications/${notificationId}/read`).then((res) => res.data)
  },
}

/**
 * 上传相关 API
 */
export const uploadApi = {
  /**
   * 上传文件
   */
  upload(file: File): Promise<ApiResponse<{ url: string }>> {
    const formData = new FormData()
    formData.append('file', file)
    return request.post('/b/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }).then((res) => res.data)
  },
}

/**
 * 菜单管理相关 API
 */
export const menuApi = {
  /**
   * 获取菜单树（管理用）
   */
  getMenuTree(): Promise<ApiResponse<Menu[]>> {
    return request.get('/b/menus').then((res) => res.data)
  },

  /**
   * 获取当前用户可见菜单
   */
  getUserMenus(): Promise<ApiResponse<Menu[]>> {
    return request.get('/b/menus/user').then((res) => res.data)
  },

  /**
   * 获取菜单详情
   */
  getMenu(menuId: number): Promise<ApiResponse<Menu>> {
    return request.get(`/b/menus/${menuId}`).then((res) => res.data)
  },

  /**
   * 创建菜单
   */
  createMenu(data: MenuRequest): Promise<ApiResponse<Menu>> {
    return request.post('/b/menus', data).then((res) => res.data)
  },

  /**
   * 更新菜单
   */
  updateMenu(menuId: number, data: MenuRequest): Promise<ApiResponse<null>> {
    return request.put(`/b/menus/${menuId}`, data).then((res) => res.data)
  },

  /**
   * 删除菜单
   */
  deleteMenu(menuId: number): Promise<ApiResponse<null>> {
    return request.delete(`/b/menus/${menuId}`).then((res) => res.data)
  },

  /**
   * 批量更新菜单排序
   */
  batchUpdateSort(data: BatchUpdateSortRequest): Promise<ApiResponse<null>> {
    return request.put('/b/menus/sort', data).then((res) => res.data)
  },
}

/**
 * 轮播图管理相关 API
 */
export const bannerApi = {
  /**
   * 获取轮播图列表
   */
  getBanners(params?: PaginationParams): Promise<ApiResponse<PaginationData<Banner>>> {
    return request.get('/b/banners', { params }).then((res) => res.data)
  },

  /**
   * 获取轮播图详情
   */
  getBanner(bannerId: number): Promise<ApiResponse<Banner>> {
    return request.get(`/b/banners/${bannerId}`).then((res) => res.data)
  },

  /**
   * 创建轮播图
   */
  createBanner(data: BannerRequest): Promise<ApiResponse<Banner>> {
    return request.post('/b/banners', data).then((res) => res.data)
  },

  /**
   * 更新轮播图
   */
  updateBanner(bannerId: number, data: BannerRequest): Promise<ApiResponse<Banner>> {
    return request.put(`/b/banners/${bannerId}`, data).then((res) => res.data)
  },

  /**
   * 删除轮播图
   */
  deleteBanner(bannerId: number): Promise<ApiResponse<null>> {
    return request.delete(`/b/banners/${bannerId}`).then((res) => res.data)
  },
}

/**
 * 统计相关 API
 */
export const statisticsApi = {
  getDashboardStats(): Promise<ApiResponse<DashboardStats>> {
    return request.get('/b/statistics/dashboard').then((res) => res.data)
  },

  getTaskStats(stationId?: number): Promise<ApiResponse<TaskStatsData>> {
    return request.get('/b/statistics/tasks', { params: { station_id: stationId } }).then((res) => res.data)
  },

  getRequestStats(stationId?: number): Promise<ApiResponse<RequestStatsData>> {
    return request.get('/b/statistics/requests', { params: { station_id: stationId } }).then((res) => res.data)
  },

  getTodayStats(): Promise<ApiResponse<TodayStatsData>> {
    return request.get('/b/statistics/today').then((res) => res.data)
  },

  getOverviewStats(params?: { days?: number; station_id?: number }): Promise<ApiResponse<OverviewStatsData>> {
    return request.get('/b/statistics/overview', { params }).then((res) => res.data)
  },

  getServiceTypeStats(params?: { days?: number; station_id?: number }): Promise<ApiResponse<ServiceTypeStatsData[]>> {
    return request.get('/b/statistics/service-types', { params }).then((res) => res.data)
  },

  getRequestTrend(params?: { days?: number; station_id?: number }): Promise<ApiResponse<TrendItemData[]>> {
    return request.get('/b/statistics/trend', { params }).then((res) => res.data)
  },

  getEfficiencyStats(params?: { days?: number; station_id?: number }): Promise<ApiResponse<EfficiencyStatsData>> {
    return request.get('/b/statistics/efficiency', { params }).then((res) => res.data)
  },

  getStaffRanking(params?: { days?: number; limit?: number; station_id?: number }): Promise<ApiResponse<StaffRankingItemData[]>> {
    return request.get('/b/statistics/staff-ranking', { params }).then((res) => res.data)
  },
}

export const reportApi = {
  generateReport(data: GenerateReportRequest): Promise<Blob> {
    return request.post('/b/reports/generate', data, { responseType: 'blob' }).then((res) => res.data)
  },

  getReports(params?: { page?: number; page_size?: number; type?: string }): Promise<ApiResponse<PaginationData<ReportData>>> {
    return request.get('/b/reports', { params }).then((res) => res.data)
  },

  downloadReport(id: number): Promise<Blob> {
    return request.get(`/b/reports/${id}/download`, { responseType: 'blob' }).then((res) => res.data)
  },

  deleteReport(id: number): Promise<ApiResponse<null>> {
    return request.delete(`/b/reports/${id}`).then((res) => res.data)
  },
}

/**
 * 新闻管理相关 API
 */
export const newsApi = {
  /**
   * 获取新闻列表
   */
  getNewsList(params?: NewsListParams): Promise<ApiResponse<PaginationData<News>>> {
    return request.get('/b/news', { params }).then((res) => res.data)
  },

  /**
   * 获取新闻详情
   */
  getNews(newsId: number): Promise<ApiResponse<News>> {
    return request.get(`/b/news/${newsId}`).then((res) => res.data)
  },

  /**
   * 创建新闻
   */
  createNews(data: NewsRequest): Promise<ApiResponse<News>> {
    return request.post('/b/news', data).then((res) => res.data)
  },

  /**
   * 更新新闻
   */
  updateNews(newsId: number, data: NewsRequest): Promise<ApiResponse<News>> {
    return request.put(`/b/news/${newsId}`, data).then((res) => res.data)
  },

  /**
   * 删除新闻
   */
  deleteNews(newsId: number): Promise<ApiResponse<null>> {
    return request.delete(`/b/news/${newsId}`).then((res) => res.data)
  },
}
