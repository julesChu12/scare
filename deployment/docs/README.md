# 部署文档

**部署方式**：GitHub Actions CI/CD + Docker Compose

---

## 架构概览

```
GitHub (master push/merge)
  → GitHub Actions (构建前端)
    → rsync 到服务器
      → Docker Compose (后端构建 + 启动全部服务)

服务器组件:
  ┌─ Nginx (:80)  ─→  静态文件 (C端 + 管理门户)
  │                ─→  反向代理 /api/ → Backend
  ├─ Backend (:8080)
  ├─ MySQL (:3306, 仅本地)
  └─ Redis (:6379, 仅本地)
```

---

## 首次部署

### 1. 服务器准备（Debian 12）

```bash
# SSH 到服务器，下载并运行初始化脚本
# 或手动创建目录 + .env 文件
curl -sL https://raw.githubusercontent.com/<你的用户名>/<仓库名>/master/deployment/server-init.sh | bash

# 编辑环境变量
nano /opt/scare/deployment/.env
```

### 2. 配置 GitHub Secrets

在 GitHub 仓库 → Settings → Secrets and variables → Actions 中添加：

| Secret 名称 | 说明 | 示例 |
|-------------|------|------|
| `SERVER_HOST` | 服务器 IP | `123.45.67.89` |
| `SERVER_USER` | SSH 用户名 | `root` |
| `SERVER_SSH_KEY` | SSH 私钥（完整内容） | `-----BEGIN OPENSSH PRIVATE KEY-----...` |
| `SERVER_PORT` | SSH 端口（可选，默认 22） | `22` |
| `AMAP_KEY` | 高德地图 Key（前端构建用，可选） | `your_amap_key` |

**SSH 密钥生成**（本地执行）：

```bash
ssh-keygen -t ed25519 -f ~/.ssh/scare_deploy -C "scare-deploy"

# 公钥添加到服务器
ssh-copy-id -i ~/.ssh/scare_deploy.pub root@你的服务器IP

# 私钥内容复制到 GitHub Secrets 的 SERVER_SSH_KEY
cat ~/.ssh/scare_deploy
```

### 3. 触发部署

推送到 master 分支即可自动触发：

```bash
git push origin master
```

---

## 日常运维

### 查看服务状态

```bash
cd /opt/scare/deployment
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.yml logs -f backend
```

### 手动重启

```bash
cd /opt/scare/deployment
docker-compose -f docker-compose.prod.yml restart backend
```

### 查看部署日志

GitHub 仓库 → Actions 标签页查看每次部署的完整日志。

---

## 文件说明

| 文件 | 说明 |
|------|------|
| `docker-compose.prod.yml` | 生产环境 Docker Compose（v2.4 格式，兼容 docker-compose 1.29） |
| `Dockerfile.backend` | 后端多阶段构建镜像 |
| `nginx/default.conf` | Nginx 反向代理 + 静态文件配置 |
| `server-init.sh` | 服务器首次初始化脚本 |
| `server-deploy.sh` | CI/CD 调用的部署脚本 |
| `.env.example` | 环境变量模板 |
| `dist/` | 前端构建产物（由 CI 自动同步，不提交到 Git） |

---

## CI/CD 流水线详情

### 触发条件

- `master` 分支的 `push` 事件（包括 merge PR）

### 流水线步骤

1. **build-frontend**：在 GitHub runner 上构建 C 端和管理门户
2. **deploy**：
   - 下载前端构建产物
   - rsync 同步项目文件到服务器 `/opt/scare/`
   - SSH 执行 `deployment/server-deploy.sh`
   - 健康检查确认部署成功

### 并发控制

同一时间只允许一个部署运行，新触发的会取消进行中的。

---

**最后更新**：2026-03-09
