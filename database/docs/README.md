# 数据库文档

**数据库**：MySQL 8.0  
**设计模式**：Code First (GORM AutoMigrate)

---

## 📊 数据库设计

详细设计文档：[backend/docs/03-数据库设计.md](../../backend/docs/03-数据库设计.md)

---

## 📁 SQL 脚本

### 表结构定义
📍 [schema/schema.sql](../schema/schema.sql)

包含所有表结构定义，由GORM自动生成

### 测试数据
📍 [seeds/seed.sql](../seeds/seed.sql)

包含测试数据：
- 12个测试用户
- 3个服务站点
- 5个地理围栏
- 角色权限配置

---

## 🗄️ 数据库表

| 表名 | 说明 |
|------|------|
| users | 用户表 |
| user_roles | 用户角色关联表 |
| customer_profiles | 客户档案表 |
| service_stations | 服务站点表 |
| service_zones | 地理围栏表 |
| service_requests | 服务需求表 |
| task_assignments | 任务分配表 |
| task_histories | 任务历史表 |
| notifications | 通知表 |
| news | 新闻表 |

---

## 🚀 快速开始

```bash
# 使用Docker Compose自动初始化
docker-compose up -d

# 手动连接数据库
docker exec -it scare_mysql mysql -uscare_user -pscare_pass scare_db

# 查看表结构
SHOW TABLES;
DESC users;

# 查看测试数据
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM service_stations;
```

---

**最后更新**：2026-01-31
