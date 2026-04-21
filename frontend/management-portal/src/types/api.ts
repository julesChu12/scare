/**
 * API 响应基础结构
 */
export interface ApiResponse<T = any> {
  msg: string
  data: T
}

/**
 * 分页参数
 */
export interface PaginationParams {
  page?: number
  page_size?: number
  station_id?: number
}

export interface StationListParams extends PaginationParams {
  keyword?: string
  name?: string
  status?: string
}

/**
 * 任务列表查询参数
 */
export interface TaskListParams extends PaginationParams {
  status?: string
  station_id?: number
  service_type?: string
  staff_id?: number
  request_no?: string
}

/**
 * 分页响应数据
 */
export interface PaginationData<T> {
  items: T[]
  page: number
  page_size: number
  total: number
}

/**
 * B端用户信息
 */
export interface User {
  id: number
  phone: string
  name: string
  email: string
  avatar?: string
  gender?: string
  birth_date?: string
  id_card?: string
  // old role fields deprecated
  // roles: string[]
  // primary_role: string

  // new identity fields
  primary_identity: string
  b_end_identities: string[]
  c_end_identities: string[]

  station_name?: string
  station_id: number | null
  status: string
  created_at: string
  updated_at: string
}

/**
 * 登录请求
 */
export interface LoginRequest {
  phone: string
  password: string
}

/**
 * 登录响应
 */
export interface LoginResponse {
  token: string
  refresh_token: string
  user_id: number
  roles: string[]
  identities?: string[] // Backend returns identities instead of roles in some cases
  type: string
  station_id: number | null
  name: string
  phone: string
  status: string
}

/**
 * 刷新 Token 响应
 */
export interface RefreshResponse {
  token: string
  refresh_token: string
}

/**
 * 当前用户信息响应
 */
export interface MeResponse {
  user: {
    id: number
    name: string
    phone: string
    roles: string[]
    identities?: string[] // Backend returns identities
    station_id: number | null
    status: string
  }
  permissions: string[]
}

/**
 * 任务信息
 */
export interface Task {
  id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  request_id: number
  station_id: number
  staff_id: number | null
  status: string
  claimed_at: string | null
  completed_at: string | null
  rating: number
  feedback: string
  staff_notes: string
  images: string[] | null
  // 关联的服务需求信息
  request?: ServiceRequest
}

export interface ServiceRequest {
  id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  request_no: string
  user_id: number
  service_type: string
  status: string
  priority?: string
  urgency?: string
  description: string
  submit_location_lat: number
  submit_location_lng: number
  contact_name: string
  contact_phone: string
  address: string
  scheduled_at?: string | null
  appointment_time?: string | null
  station_id: number | null
  station_name?: string
}

export interface UpdateServiceRequest {
  service_type?: string
  description?: string
  contact_name?: string
  contact_phone?: string
  address?: string
  priority?: string
  urgency?: string
  scheduled_at?: string
  appointment_time?: string
  station_id?: number
  status?: string
}

export interface CompleteTaskRequest {
  images?: string[]
  staff_notes?: string
}

export interface Station {
  id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  name: string
  code: string
  address: string
  phone: string
  latitude: number
  longitude: number
  service_area: string
  work_hours: string
  status: string
}

export interface StationRequest {
  name: string
  code?: string
  address?: string
  phone?: string
  latitude?: number
  longitude?: number
  service_area?: string
  work_hours?: string
  status?: string
}

export interface ZonePoint {
  lat: number
  lng: number
}

export interface Zone {
  id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  station_id: number
  name: string
  points: string | ZonePoint[] // API returns JSON string, but we want to work with ZonePoint[]
  priority: number
  status: string
}

export interface ZoneRequest {
  station_id: number
  name: string
  points: ZonePoint[]
  priority?: number
  status?: string
}

export interface UserCreateRequest {
  phone: string
  password: string
  name?: string
  email?: string
  identity_type: string
  station_id?: number | null
  status?: string
}

export interface UserUpdateRequest {
  name?: string
  email?: string
  avatar?: string
  gender?: string
  birth_date?: string
  id_card?: string
  role?: string
  station_id?: number | null
  status?: string
  password?: string
}

export interface UpdateUserIdentitiesRequest {
  identities: string[]
}

export interface UpdateRolePermissionsRequest {
  permissions: string[]
}

export interface PermissionNode {
  id: string           // 权限码
  code: string         // 权限码（同 id）
  label: string        // 显示名称
  type: 'menu' | 'button' | 'resource'
  api_path?: string
  method?: string
  is_public?: boolean
  disabled?: boolean
  children?: PermissionNode[]
}

export interface Notification {
  id: number
  created_at: string
  updated_at: string
  user_id: number
  title: string
  content?: string
  body?: string
  type?: string
  is_read: boolean
  read_at: string | null
}

export interface NotificationListParams extends PaginationParams {
  type?: string
  is_read?: boolean
}

export interface Menu {
  id: number
  parent_id: number
  name: string
  path: string
  component: string
  icon: string
  permission_code: string  // 新字段名
  sort: number
  hidden: boolean
  status: string
  children?: Menu[]
  created_at?: string
  updated_at?: string
}

export interface MenuRequest {
  parent_id?: number
  name: string
  path?: string
  component?: string
  icon?: string
  permission_code?: string  // 新字段名
  sort?: number
  hidden?: boolean
  status?: string
}

export interface BatchUpdateSortRequest {
  updates: Array<{
    id: number
    sort: number
  }>
}

export interface Banner {
  id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  title: string
  image_url: string
  link_type: string
  link_value: string
  sort: number
  status: string
  station_id: number
}

/**
 * 创建/更新轮播图请求
 */
export interface BannerRequest {
  title: string
  image_url: string
  link_type: string
  link_value: string
  sort: number
  status: string
  station_id: number
}

/**
 * 任务统计数据
 */
export interface TaskStatsData {
  total: number
  dispatched: number
  claimed: number
  completed: number
  cancelled: number
}

/**
 * 需求统计数据
 */
export interface RequestStatsData {
  total: number
  pending: number
  dispatched: number
  processing: number
  completed: number
  cancelled: number
}

/**
 * 今日统计数据
 */
export interface TodayStatsData {
  new_requests: number
  completed_tasks: number
  new_users: number
  avg_response_time: number
}

/**
 * 我的任务统计
 */
export interface MyTaskStatsData {
  claimed: number
  completed: number
  total: number
}

/**
 * 工作台统计数据
 */
export interface DashboardStats {
  task_stats: TaskStatsData
  request_stats: RequestStatsData
  today_stats: TodayStatsData
  my_task_stats?: MyTaskStatsData
}

/**
 * 新闻类型
 */
export type NewsType = 'news' | 'notice' | 'activity'

/**
 * 新闻状态
 */
export type NewsStatus = 'draft' | 'published' | 'archived'

/**
 * 新闻信息
 */
export interface News {
  id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  title: string
  summary: string
  content: string
  cover_url: string
  type: NewsType
  status: NewsStatus
  station_id: number
  author_id: number
  publish_at: string
  view_count: number
}

/**
 * 创建/更新新闻请求
 */
export interface NewsRequest {
  title: string
  summary?: string
  content?: string
  cover_url?: string
  type?: NewsType
  status?: NewsStatus
  station_id?: number
}

/**
 * 新闻列表查询参数
 */
export interface NewsListParams extends PaginationParams {
  type?: NewsType
  status?: NewsStatus
  station_id?: number
}

export interface OverviewStatsData {
  total_requests: number
  pending: number
  completed: number
  in_progress: number
}

export interface ServiceTypeStatsData {
  type: string
  name: string
  count: number
  percentage: number
}

export interface TrendItemData {
  date: string
  count: number
  percentage: number
}

export interface EfficiencyStatsData {
  avg_response_time: number
  avg_process_time: number
  satisfaction_rate: number
}

export interface StaffRankingItemData {
  id: number
  name: string
  completed_count: number
  avg_rating: number
  is_online: boolean
}

export interface GenerateReportRequest {
  type: 'service' | 'performance' | 'request' | 'station'
  format: 'xlsx' | 'csv'
  station_id?: number | null
  start_date: string
  end_date: string
  preview?: boolean
}

export interface ReportListParams extends PaginationParams {
  type?: string
}

export interface ReportData {
  id: number
  created_at: string
  name: string
  type: string
  format: string
  file_size: number
  station_id: number | null
  start_date: string
  end_date: string
  created_by: number
}

export interface ReportPreviewData {
  request_count: number
  completed_request_count: number
  completion_rate: number
  service_type_count: number
  ranked_staff_count: number
  completed_task_count: number
  avg_rating: number
  trend_days: number
  station_count: number
}

/**
 * 老年人档案信息
 */
export interface ElderlyProfile {
  id: number
  name: string
  phone: string
  gender: string
  birth_date: string
  id_card: string
  address: string
  station_id: number
  station_name: string
  health_status: string
  disability_level: string
  medical_history: string
  special_needs: string
  customer_type: string
  created_at: string
}

/**
 * 创建老人档案请求
 */
export interface ElderlyCreateRequest {
  name: string
  phone: string
  gender?: string
  birth_date?: string
  id_card?: string
  address?: string
  station_id?: number
  health_status?: string
  disability_level?: string
  medical_history?: string
  special_needs?: string
}

/**
 * 更新老人档案请求
 */
export interface ElderlyUpdateRequest {
  name?: string
  gender?: string
  birth_date?: string
  id_card?: string
  address?: string
  station_id?: number
  health_status?: string
  disability_level?: string
  medical_history?: string
  special_needs?: string
}

/**
 * 老人档案列表查询参数
 */
export interface ElderlyListParams extends PaginationParams {
  keyword?: string
  station_id?: number
  health_status?: string
}

/**
 * 老人服务记录
 */
export interface ElderlyServiceRecord {
  request_id: number
  request_no: string
  service_type: string
  status: string
  description: string
  address: string
  rating: number
  feedback: string
  created_at: string
  task_id: number | null
  staff_name: string | null
  claimed_at: string | null
  completed_at: string | null
}
