## Admin 登录
curl -s -X "POST" "http://localhost:8080/api/v1/b/auth/login" \
     -H "Content-Type: application/json" \
     -d '{"phone":"13800000001","password":"Test@123"}'

## Station Manager 登录
curl -s -X "POST" "http://localhost:8080/api/v1/b/auth/login" \
     -H "Content-Type: application/json" \
     -d '{"phone":"13800000002","password":"Test@123"}'

## Staff 登录
curl -s -X "POST" "http://localhost:8080/api/v1/b/auth/login" \
     -H "Content-Type: application/json" \
     -d '{"phone":"13800000004","password":"Test@123"}'

## C端老年人登录
curl -s -X "POST" "http://localhost:8080/api/v1/c/auth/login" \
     -H "Content-Type: application/json" \
     -d '{"phone":"13800000008","password":"Test@123"}'

## 健康检查（公开）
curl -s "http://localhost:8080/api/v1/health"

## 获取当前用户信息
curl -s "http://localhost:8080/api/v1/b/auth/me" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 登出
curl -s -X "POST" "http://localhost:8080/api/v1/b/auth/logout" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 刷新 Token
curl -s -X "POST" "http://localhost:8080/api/v1/b/auth/refresh" \
     -H "Content-Type: application/json" \
     -d '{"refresh_token":"<paste_refresh_token>"}'

## 用户列表
curl -s "http://localhost:8080/api/v1/b/users" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 用户列表（分页）
curl -s "http://localhost:8080/api/v1/b/users?page=1&page_size=10" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 用户列表（按角色筛选）
curl -s "http://localhost:8080/api/v1/b/users?role=staff" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 用户列表（按站点筛选）
curl -s "http://localhost:8080/api/v1/b/users?station_id=1" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 用户详情
curl -s "http://localhost:8080/api/v1/b/users/1" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 创建用户
curl -s -X "POST" "http://localhost:8080/api/v1/b/users" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "phone": "13800000099",
    "password": "Test@123",
    "name": "测试用户",
    "role": "staff",
    "station_id": 1
  }'

## 更新用户
curl -s -X "PUT" "http://localhost:8080/api/v1/b/users/1" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "name": "系统管理员",
    "email": "admin@scare.com"
  }'

## 权限树
curl -s "http://localhost:8080/api/v1/b/permissions/tree" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 获取角色权限（admin）
curl -s "http://localhost:8080/api/v1/b/roles/admin/permissions" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 获取角色权限（station_manager）
curl -s "http://localhost:8080/api/v1/b/roles/station_manager/permissions" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 获取角色权限（staff）
curl -s "http://localhost:8080/api/v1/b/roles/staff/permissions" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 更新角色权限
curl -s -X "PUT" "http://localhost:8080/api/v1/b/roles/staff/permissions" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "permission_ids": [10, 11, 12, 13, 14, 15, 16, 17, 18, 19]
  }'

## 更新用户身份
curl -s -X "PUT" "http://localhost:8080/api/v1/b/users/4/identities" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "identities": [
      {"identity_type": "staff", "is_primary": true, "station_id": 1}
    ]
  }'

## 获取用户菜单
curl -s "http://localhost:8080/api/v1/b/menus/user" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 菜单列表（树形）
curl -s "http://localhost:8080/api/v1/b/menus" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 菜单详情
curl -s "http://localhost:8080/api/v1/b/menus/1" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 创建菜单
curl -s -X "POST" "http://localhost:8080/api/v1/b/menus" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "name": "测试菜单",
    "path": "/test",
    "component": "test/index",
    "parent_id": 0,
    "sort": 99,
    "icon": "Setting",
    "permission_code": "test"
  }'

## 更新菜单
curl -s -X "PUT" "http://localhost:8080/api/v1/b/menus/1" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "name": "工作台",
    "sort": 1
  }'

## 删除菜单
curl -s -X "DELETE" "http://localhost:8080/api/v1/b/menus/999" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 批量排序
curl -s -X "PUT" "http://localhost:8080/api/v1/b/menus/sort" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "items": [{"id": 1, "sort": 1}, {"id": 2, "sort": 2}]
  }'

## 站点列表
curl -s "http://localhost:8080/api/v1/b/stations" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 站点详情
curl -s "http://localhost:8080/api/v1/b/stations/1" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 创建站点
curl -s -X "POST" "http://localhost:8080/api/v1/b/stations" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "name": "测试站点",
    "code": "TEST-001",
    "address": "北京市昌平区测试地址",
    "phone": "010-12345678",
    "latitude": 40.05,
    "longitude": 116.38,
    "capacity": 10,
    "work_hours": "08:00-18:00"
  }'

## 更新站点
curl -s -X "PUT" "http://localhost:8080/api/v1/b/stations/1" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "name": "霍营街道养老服务站",
    "capacity": 15
  }'

## 删除站点
curl -s -X "DELETE" "http://localhost:8080/api/v1/b/stations/999" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 围栏列表
curl -s "http://localhost:8080/api/v1/b/zones" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 围栏列表（按站点筛选）
curl -s "http://localhost:8080/api/v1/b/zones?station_id=1" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 创建围栏
curl -s -X "POST" "http://localhost:8080/api/v1/b/zones" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "station_id": 1,
    "name": "测试围栏",
    "points": [[116.37, 40.04], [116.39, 40.04], [116.39, 40.06], [116.37, 40.06]],
    "priority": 1
  }'

## 更新围栏
curl -s -X "PUT" "http://localhost:8080/api/v1/b/zones/1" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "name": "霍营核心区域",
    "priority": 10
  }'

## 删除围栏
curl -s -X "DELETE" "http://localhost:8080/api/v1/b/zones/999" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## B端 - 服务请求管理

## 服务请求列表
curl -s "http://localhost:8080/api/v1/b/requests" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 服务请求列表（分页+筛选）
curl -s "http://localhost:8080/api/v1/b/requests?page=1&page_size=10&status=pending" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 服务请求详情
curl -s "http://localhost:8080/api/v1/b/requests/1" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 更新请求状态（派单）
curl -s -X "PUT" "http://localhost:8080/api/v1/b/requests/4/status" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{"status": "dispatched", "station_id": 1}'

## 任务列表
curl -s "http://localhost:8080/api/v1/b/tasks" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 任务池
curl -s "http://localhost:8080/api/v1/b/tasks/pool" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 任务池（所有站点）
curl -s "http://localhost:8080/api/v1/b/tasks/pool?station_id=0" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 任务池（指定站点）
curl -s "http://localhost:8080/api/v1/b/tasks/pool?station_id=1" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 我的任务
curl -s "http://localhost:8080/api/v1/b/tasks/my" \
     -H "Authorization: Bearer <STAFF_TOKEN>"

## 任务详情
curl -s "http://localhost:8080/api/v1/b/tasks/1" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 认领任务
curl -s -X "POST" "http://localhost:8080/api/v1/b/tasks/3/claim" \
     -H "Authorization: Bearer <STAFF_TOKEN>"

## 完成任务
curl -s -X "POST" "http://localhost:8080/api/v1/b/tasks/1/complete" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <STAFF_TOKEN>" \
     -d '{"staff_notes": "任务已完成"}'

## 转移任务
curl -s -X "POST" "http://localhost:8080/api/v1/b/tasks/2/transfer" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <SM_TOKEN>" \
     -d '{"to_staff_id": 4, "remark": "转给王小红处理"}'

## 通知列表
curl -s "http://localhost:8080/api/v1/b/notifications" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 通知列表（分页）
curl -s "http://localhost:8080/api/v1/b/notifications?page=1&page_size=10" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 标记通知已读
curl -s -X "POST" "http://localhost:8080/api/v1/b/notifications/1/read" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## B端 - 轮播图管理

## 轮播图列表
curl -s "http://localhost:8080/api/v1/b/banners" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 创建轮播图
curl -s -X "POST" "http://localhost:8080/api/v1/b/banners" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "title": "测试轮播图",
    "image_url": "/static/test.jpg",
    "link_type": "none",
    "sort": 1,
    "station_id": 0
  }'

## 更新轮播图
curl -s -X "PUT" "http://localhost:8080/api/v1/b/banners/1" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{"title": "更新标题", "sort": 2}'

## 删除轮播图
curl -s -X "DELETE" "http://localhost:8080/api/v1/b/banners/999" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## B端 - 新闻管理

## 新闻列表
curl -s "http://localhost:8080/api/v1/b/news" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 新闻详情
curl -s "http://localhost:8080/api/v1/b/news/1" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 创建新闻
curl -s -X "POST" "http://localhost:8080/api/v1/b/news" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "title": "测试新闻",
    "summary": "这是一条测试新闻",
    "content": "<p>新闻正文内容</p>",
    "type": "news",
    "status": "published",
    "station_id": 0
  }'

## 更新新闻
curl -s -X "PUT" "http://localhost:8080/api/v1/b/news/1" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{"title": "更新标题"}'

## 删除新闻
curl -s -X "DELETE" "http://localhost:8080/api/v1/b/news/999" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## B端 - 统计数据

## 仪表盘统计
curl -s "http://localhost:8080/api/v1/b/statistics/dashboard" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 任务统计
curl -s "http://localhost:8080/api/v1/b/statistics/tasks" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 请求统计
curl -s "http://localhost:8080/api/v1/b/statistics/requests" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 今日统计
curl -s "http://localhost:8080/api/v1/b/statistics/today" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 概览统计
curl -s "http://localhost:8080/api/v1/b/statistics/overview" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 服务类型统计
curl -s "http://localhost:8080/api/v1/b/statistics/service-types" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 请求趋势
curl -s "http://localhost:8080/api/v1/b/statistics/trend" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 效率统计
curl -s "http://localhost:8080/api/v1/b/statistics/efficiency" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 员工排名
curl -s "http://localhost:8080/api/v1/b/statistics/staff-ranking" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## B端 - 报表

## 生成报表
curl -s -X "POST" "http://localhost:8080/api/v1/b/reports/generate" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -d '{
    "type": "service",
    "format": "xlsx",
    "start_date": "2026-01-01",
    "end_date": "2026-01-31"
  }'

## 报表列表
curl -s "http://localhost:8080/api/v1/b/reports" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 下载报表
curl -s "http://localhost:8080/api/v1/b/reports/1/download" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## 删除报表
curl -s -X "DELETE" "http://localhost:8080/api/v1/b/reports/999" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## B端 - 文件上传

## 上传文件（multipart/form-data）
curl -s -X "POST" "http://localhost:8080/api/v1/b/upload" \
     -H "Authorization: Bearer <ADMIN_TOKEN>" \
     -F "file=@/path/to/your/file.jpg"

## C端 - 认证（公开）

## 发送验证码
curl -s -X "POST" "http://localhost:8080/api/v1/c/auth/send-code" \
     -H "Content-Type: application/json" \
     -d '{"phone":"13800000008"}'

## 快速开始
curl -s -X "POST" "http://localhost:8080/api/v1/c/auth/quick-start" \
     -H "Content-Type: application/json" \
     -d '{"phone":"13800000008"}'

## C端登录
curl -s -X "POST" "http://localhost:8080/api/v1/c/auth/login" \
     -H "Content-Type: application/json" \
     -d '{"phone":"13800000008","password":"Test@123"}'

## C端刷新 Token
curl -s -X "POST" "http://localhost:8080/api/v1/c/auth/refresh" \
     -H "Content-Type: application/json" \
     -d '{"refresh_token":"<paste_c_refresh_token>"}'

## C端 - 公开接口（无需认证）

## 新闻列表
curl -s "http://localhost:8080/api/v1/c/news"

## 新闻详情
curl -s "http://localhost:8080/api/v1/c/news/1"

## 轮播图列表
curl -s "http://localhost:8080/api/v1/c/banners"

## 站点匹配（GET）
curl -s "http://localhost:8080/api/v1/c/stations/match?lng=116.38&lat=40.05"

## 站点匹配（POST）
curl -s -X "POST" "http://localhost:8080/api/v1/c/stations/match" \
     -H "Content-Type: application/json" \
     -d '{"lng": 116.38, "lat": 40.05}'

## 地理编码
curl -s -X "POST" "http://localhost:8080/api/v1/c/geocode" \
     -H "Content-Type: application/json" \
     -d '{"address": "北京市昌平区霍营街道"}'

## 逆地理编码（GET）
curl -s "http://localhost:8080/api/v1/c/geocode/reverse?lng=116.38&lat=40.05"

## 逆地理编码（POST）
curl -s -X "POST" "http://localhost:8080/api/v1/c/geocode/reverse" \
     -H "Content-Type: application/json" \
     -d '{"lng": 116.38, "lat": 40.05}'

## C端 - 需认证接口

## 获取当前用户
curl -s "http://localhost:8080/api/v1/c/auth/me" \
     -H "Authorization: Bearer <C_TOKEN>"

## 检查 Token
curl -s "http://localhost:8080/api/v1/c/auth/check" \
     -H "Authorization: Bearer <C_TOKEN>"

## C端登出
curl -s -X "POST" "http://localhost:8080/api/v1/c/auth/logout" \
     -H "Authorization: Bearer <C_TOKEN>"

## 更新个人资料
curl -s -X "PUT" "http://localhost:8080/api/v1/c/profile" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <C_TOKEN>" \
     -d '{
    "name": "张大爷",
    "address": "北京市昌平区霍营街道华龙苑北里小区",
    "latitude": 40.05,
    "longitude": 116.38
  }'

## 创建服务请求
curl -s -X "POST" "http://localhost:8080/api/v1/c/requests" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <C_TOKEN>" \
     -d '{
    "service_type": "meal",
    "description": "需要午餐送餐服务",
    "contact_name": "张大爷",
    "contact_phone": "13800000008",
    "address": "北京市昌平区霍营街道华龙苑北里小区1号楼",
    "submit_location_lat": 40.05,
    "submit_location_lng": 116.38,
    "urgency": "normal"
  }'

## 我的服务请求列表
curl -s "http://localhost:8080/api/v1/c/requests" \
     -H "Authorization: Bearer <C_TOKEN>"

## 服务请求详情
curl -s "http://localhost:8080/api/v1/c/requests/1" \
     -H "Authorization: Bearer <C_TOKEN>"

## 取消服务请求
curl -s -X "POST" "http://localhost:8080/api/v1/c/requests/4/cancel" \
     -H "Authorization: Bearer <C_TOKEN>"

## 评价服务
curl -s -X "POST" "http://localhost:8080/api/v1/c/requests/1/rate" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <C_TOKEN>" \
     -d '{"rating": 5, "feedback": "服务很好"}'

## C端通知列表
curl -s "http://localhost:8080/api/v1/c/notifications" \
     -H "Authorization: Bearer <C_TOKEN>"

## C端标记通知已读
curl -s -X "POST" "http://localhost:8080/api/v1/c/notifications/1/read" \
     -H "Authorization: Bearer <C_TOKEN>"

## C端文件上传
curl -s -X "POST" "http://localhost:8080/api/v1/c/upload" \
     -H "Authorization: Bearer <C_TOKEN>" \
     -F "file=@/path/to/your/file.jpg"

## 权限隔离验证

## 未认证访问 B端（应返回 401）
curl -s "http://localhost:8080/api/v1/b/auth/me"

## B端 Token 访问 C端（应返回 403）
curl -s "http://localhost:8080/api/v1/c/auth/me" \
     -H "Authorization: Bearer <ADMIN_TOKEN>"

## C端 Token 访问 B端（应返回 403）
curl -s "http://localhost:8080/api/v1/b/auth/me" \
     -H "Authorization: Bearer <C_TOKEN>"

## Staff 访问角色权限管理（应返回 403）
curl -s "http://localhost:8080/api/v1/b/roles/admin/permissions" \
     -H "Authorization: Bearer <STAFF_TOKEN>"
