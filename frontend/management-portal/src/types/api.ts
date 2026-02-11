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
 * B端用户角色
 */
export type UserRole = 'admin' | 'station_manager' | 'staff'

/**
 * 用户状态
 */
export type UserStatus = 'active' | 'inactive'

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
  status: UserStatus
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
 * 任务状态
 */
export type TaskStatus = 'dispatched' | 'claimed' | 'completed' | 'cancelled'

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
  status: TaskStatus
  claimed_at: string | null
  completed_at: string | null
  rating: number
  feedback: string
  staff_notes: string
  images: string[] | null
  // 关联的服务需求信息
  request?: ServiceRequest
}

export type RequestStatus = 'pending' | 'dispatched' | 'claimed' | 'processing' | 'completed' | 'cancelled' | 'rejected'

export interface ServiceRequest {
  id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
  request_no: string
  user_id: number
  service_type: string
  status: RequestStatus
  priority: string
  description: string
  submit_location_lat: number
  submit_location_lng: number
  contact_name: string
  contact_phone: string
  address: string
  scheduled_at: string | null
  station_id: number | null
  station_name?: string
}

export interface CreateServiceRequest {
  user_id?: number
  service_type: string
  description: string
  contact_name: string
  contact_phone: string
  address: string
  priority?: string
  scheduled_at?: string
  station_id?: number
}

export interface UpdateServiceRequest {
  service_type?: string
  description?: string
  contact_name?: string
  contact_phone?: string
  address?: string
  priority?: string
  scheduled_at?: string
  station_id?: number
  status?: RequestStatus
}

export interface CompleteTaskRequest {
  images?: string[]
  staff_notes?: string
}

export type StationStatus = 'active' | 'inactive'

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
  status: StationStatus
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
  role: string
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

export interface UpdateUserRolesRequest {
  roles: string[]
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
  content: string
  type: string
  is_read: boolean
  read_at: string | null
}

export type MenuStatus = 'active' | 'inactive'

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
  status: MenuStatus
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
  status?: MenuStatus
}

export interface BatchUpdateSortRequest {
  updates: Array<{
    id: number
    sort: number
  }>
}

export type BannerStatus = 'active' | 'inactive'

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
  status: BannerStatus
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
  status: BannerStatus
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
