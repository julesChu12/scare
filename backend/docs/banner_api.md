# Banner 轮播图 API 文档

## 功能概述

Banner 轮播图系统支持：
- 全局 Banner（所有站点可见，`station_id = 0`）
- 站点专属 Banner（特定站点可见，`station_id > 0`）
- 优先级排序（站点专属 > 全局，按 `sort` 降序）
- 链接类型：无链接、外部URL、新闻详情、服务详情
- 状态管理：激活/非激活

## 数据模型

```go
type Banner struct {
    ID        int64     // 主键
    StationID int64     // 站点ID（0表示全局）
    Title     string    // 标题
    ImageURL  string    // 图片URL（必填）
    LinkType  string    // 链接类型：none, url, news, service
    LinkValue string    // 链接值
    Sort      int64     // 排序值（降序）
    Status    string    // 状态：active, inactive
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt time.Time // 软删除
}
```

## C端 API（公开接口）

### 获取轮播图列表

**接口**: `GET /api/v1/c/banners`

**参数**:
- `station_id` (可选): 站点ID，0或不传表示获取全局 Banner

**响应示例**:
```json
{
  "msg": "ok",
  "data": [
    {
      "id": 1,
      "station_id": 0,
      "title": "社区养老服务平台",
      "image_url": "https://example.com/banner1.jpg",
      "link_type": "none",
      "link_value": "",
      "sort": 100,
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**业务逻辑**:
- 如果传入 `station_id > 0`：返回该站点专属 + 全局 Banner
- 如果 `station_id = 0` 或不传：仅返回全局 Banner
- 只返回 `status = 'active'` 的 Banner
- 排序：站点专属优先，然后按 `sort` 降序

## B端 API（管理接口，需要认证）

### 1. 获取 Banner 列表（分页）

**接口**: `GET /api/v1/b/banners`

**权限**: 需要 `banner:list` 权限

**参数**:
- `page` (可选): 页码，默认 1
- `page_size` (可选): 每页数量，默认 10
- `station_id` (可选): 站点ID筛选

**响应示例**:
```json
{
  "msg": "ok",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

### 2. 创建 Banner

**接口**: `POST /api/v1/b/banners`

**权限**: 需要 `banner:create` 权限

**请求体**:
```json
{
  "station_id": 0,
  "title": "新年活动",
  "image_url": "https://example.com/banner.jpg",
  "link_type": "url",
  "link_value": "https://example.com/activity",
  "sort": 100,
  "status": "active"
}
```

**字段说明**:
- `image_url`: 必填
- `link_type`: 可选，默认 "none"，可选值：none, url, news, service
- `status`: 可选，默认 "active"

### 3. 更新 Banner

**接口**: `PUT /api/v1/b/banners/:id`

**权限**: 需要 `banner:update` 权限

**请求体**: 同创建接口

### 4. 删除 Banner

**接口**: `DELETE /api/v1/b/banners/:id`

**权限**: 需要 `banner:delete` 权限

**说明**: 软删除，数据不会真正删除

## 链接类型说明

| LinkType | 说明 | LinkValue 示例 |
|----------|------|----------------|
| `none` | 无链接 | 空字符串 |
| `url` | 外部链接 | `https://example.com` |
| `news` | 新闻详情 | 新闻ID，如 `"123"` |
| `service` | 服务详情 | 服务ID，如 `"456"` |

## 前端集成示例

### C端 - 首页轮播图

```typescript
// 获取轮播图
const fetchBanners = async (stationId?: number) => {
  const params = stationId ? { station_id: stationId } : {}
  const res = await request.get('/c/banners', { params })
  return res.data
}

// 处理点击事件
const handleBannerClick = (banner: Banner) => {
  switch (banner.link_type) {
    case 'url':
      window.open(banner.link_value, '_blank')
      break
    case 'news':
      router.push(`/news/${banner.link_value}`)
      break
    case 'service':
      router.push(`/service/${banner.link_value}`)
      break
    case 'none':
    default:
      // 无操作
      break
  }
}
```

### B端 - Banner 管理

```typescript
// 获取列表
const fetchBanners = async (page: number, pageSize: number, stationId?: number) => {
  const params = { page, page_size: pageSize }
  if (stationId !== undefined) params.station_id = stationId
  const res = await request.get('/b/banners', { params })
  return res.data
}

// 创建
const createBanner = async (data: BannerInput) => {
  const res = await request.post('/b/banners', data)
  return res.data
}

// 更新
const updateBanner = async (id: number, data: BannerInput) => {
  const res = await request.put(`/b/banners/${id}`, data)
  return res.data
}

// 删除
const deleteBanner = async (id: number) => {
  await request.delete(`/b/banners/${id}`)
}
```

## 数据库初始化

执行初始化脚本：

```bash
mysql -u root -p your_database < database/seeds/modules/60_content.sql
```

或在 MySQL 客户端中：

```sql
source /path/to/backend/database/seeds/modules/60_content.sql;
```

## 权限配置

需要在 Casbin 策略中添加以下权限：

```csv
p, admin, /api/v1/b/banners, GET
p, admin, /api/v1/b/banners, POST
p, admin, /api/v1/b/banners/:id, PUT
p, admin, /api/v1/b/banners/:id, DELETE
```

或使用权限标识：

```csv
p, admin, banner:list, *
p, admin, banner:create, *
p, admin, banner:update, *
p, admin, banner:delete, *
```

## 注意事项

1. **图片上传**: 需要先通过 `/api/v1/b/upload` 或 `/api/v1/c/upload` 上传图片，获取 URL 后再创建 Banner
2. **站点隔离**: 站点管理员只能管理自己站点的 Banner（需要在 Service 层添加权限检查）
3. **排序规则**:
   - 站点专属 Banner 优先于全局 Banner
   - 相同类型按 `sort` 降序
   - `sort` 相同按 `id` 降序
4. **缓存建议**: C端接口建议添加 Redis 缓存，减少数据库查询

## 测试

```bash
# C端 - 获取全局 Banner
curl http://localhost:8080/api/v1/c/banners

# C端 - 获取站点 Banner
curl http://localhost:8080/api/v1/c/banners?station_id=1

# B端 - 获取列表（需要 token）
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/b/banners?page=1&page_size=10

# B端 - 创建 Banner
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"station_id":0,"title":"测试","image_url":"https://example.com/test.jpg","sort":100}' \
  http://localhost:8080/api/v1/b/banners
```
