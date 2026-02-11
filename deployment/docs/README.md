# 部署文档

**部署方式**：Docker + Docker Compose

---

## 🚀 部署方案

详细部署方案：[backend/docs/07-部署方案.md](../../backend/docs/07-部署方案.md)

---

## 📦 Docker 配置

### Docker Compose
📍 [docker-compose.yml](../../docker-compose.yml)

服务包含：
- MySQL 8.0
- Redis 7.0
- (后端和前端需要单独构建)

---

## 🔧 快速部署

### 开发环境

```bash
# 启动服务
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 生产环境

```bash
# 使用生产配置
docker-compose -f docker-compose.prod.yml up -d
```

---

## 📋 环境配置

### 必需环境变量

```bash
# 数据库
DB_HOST=localhost
DB_PORT=3306
DB_USER=scare_user
DB_PASSWORD=your_password
DB_NAME=scare_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRATION=2h
```

---

**最后更新**：2026-01-31
