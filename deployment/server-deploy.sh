#!/bin/bash
# ============================================
# sCare 服务器部署脚本
# 流程: 停止旧服务 → 启动 MySQL/Redis → 执行迁移 → 启动后端 → 健康检查
# ============================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.prod.yml"
ENV_FILE="$SCRIPT_DIR/.env"
MIGRATE_SCRIPT="$PROJECT_ROOT/scripts/migrate.sh"

# 强制项目名为 scare，确保容器名一致（scare_mysql, scare_backend 等）
COMPOSE_PROJECT="scare"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()  { echo -e "${GREEN}[DEPLOY]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()  { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# docker-compose 命令带上 -p scare
docker_compose() {
    docker-compose -p "$COMPOSE_PROJECT" -f "$COMPOSE_FILE" "$@"
}

log "===== sCare 服务器部署 ====="
log "时间: $(date '+%Y-%m-%d %H:%M:%S')"

if [ ! -f "$ENV_FILE" ]; then
    err "未找到 $ENV_FILE，请先创建环境配置文件（参考 .env.example）"
fi

# 加载环境变量 (使用 set -a 将变量导出给子进程和 docker-compose)
if [ -f "$ENV_FILE" ]; then
    set -a
    source "$ENV_FILE"
    set +a
else
    err "未找到 $ENV_FILE，请先创建环境配置文件（参考 .env.example）"
fi

if ! command -v docker &> /dev/null; then
    err "未安装 Docker"
fi

if ! command -v docker-compose &> /dev/null; then
    err "未安装 docker-compose"
fi

# 检查前端构建产物
if [ ! -d "$SCRIPT_DIR/dist/c-end" ] || [ ! -d "$SCRIPT_DIR/dist/management-portal" ]; then
    err "前端构建产物不存在，请确认 CI 流程已正确执行"
fi

# 切换到 deployment 目录，确保 docker-compose 能找到 .env
cd "$SCRIPT_DIR"

# 修复 backend/database 目录权限（Docker init 脚本可能以 root 创建）
if [ -d "$PROJECT_ROOT/backend/database" ]; then
    chmod -R 755 "$PROJECT_ROOT/backend/database" 2>/dev/null || true
fi

# ============================================
# 步骤 1: 停止旧服务（强制删除残留容器）
# ============================================
log "[1/5] 停止旧服务..."
docker_compose down --remove-orphans 2>/dev/null || true
# 强制删除同名残留容器（可能由旧版 deploy 脚本遗留）
docker rm -f scare_mysql scare_redis scare_backend scare_nginx 2>/dev/null || true

# ============================================
# 步骤 2: 启动 MySQL 和 Redis（不启动 backend，等待迁移完成）
# ============================================
log "[2/5] 启动 MySQL 和 Redis..."
docker_compose up -d mysql redis

# 等待 MySQL 健康
log "  等待 MySQL 就绪..."
for i in $(seq 1 30); do
    if docker exec scare_mysql mysqladmin ping -h localhost -u root -p"${DB_ROOT_PASSWORD}" &>/dev/null 2>&1; then
        echo -e " ${GREEN}MySQL ✅${NC}"
        break
    fi
    echo -n "."
    sleep 2
    [ "$i" -eq 30 ] && err "MySQL 启动超时"
done

# 等待应用库账号可用（迁移脚本优先使用 DB_USER/DB_PASSWORD）
log "  等待数据库业务账号就绪..."
for i in $(seq 1 30); do
    if docker exec scare_mysql mysql -h127.0.0.1 -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" -e "SELECT 1" >/dev/null 2>&1; then
        echo -e " ${GREEN}DB User ✅${NC}"
        break
    fi
    echo -n "."
    sleep 2
    [ "$i" -eq 30 ] && {
        echo ""
        log "数据库业务账号检查失败详情："
        docker exec scare_mysql mysql -h127.0.0.1 -u"${DB_USER}" -p"${DB_PASSWORD}" "${DB_NAME}" -e "SELECT 1" 2>&1 || true
        err "数据库业务账号启动超时或密码不正确"
    }
done

# 等待 Redis 健康
log "  等待 Redis 就绪..."
for i in $(seq 1 15); do
    # 尝试连接 Redis，如果不带密码则自动尝试
    if docker exec scare_redis redis-cli -a "${REDIS_PASSWORD:-}" ping 2>&1 | grep -q PONG; then
        echo -e " ${GREEN}Redis ✅${NC}"
        break
    fi
    echo -n "."
    sleep 2
    [ "$i" -eq 15 ] && {
        echo ""
        log "Redis 检查失败详情："
        docker exec scare_redis redis-cli -a "${REDIS_PASSWORD:-}" ping 2>&1 || true
        err "Redis 启动超时"
    }
done

# ============================================
# 步骤 3: 执行数据库迁移
# ============================================
log "[3/5] 执行数据库迁移..."
if [ -f "$MIGRATE_SCRIPT" ]; then
    chmod +x "$MIGRATE_SCRIPT"
    bash "$MIGRATE_SCRIPT" || err "迁移执行失败"
else
    warn "迁移脚本不存在，跳过"
fi

# ============================================
# 步骤 4: 启动后端（build）
# ============================================
log "[4/5] 构建并启动后端服务..."
cd "$SCRIPT_DIR"
docker_compose up -d --build

# ============================================
# 步骤 5: 健康检查
# ============================================
log "[5/5] 健康检查..."
echo -n "  等待后端启动"
READY=false
for i in $(seq 1 30); do
    if curl -sf http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        echo -e " ${GREEN}✅${NC}"
        READY=true
        break
    fi
    echo -n "."
    sleep 2
done

if [ "$READY" = false ]; then
    warn "后端健康检查超时，查看日志："
    docker_compose logs --tail=50 backend
    exit 1
fi

log "===== 部署完成 ====="
log "  后端 API:  http://localhost:8080"
log "  Nginx:     http://localhost"
log "  管理后台:  http://localhost/manage/"
echo ""
docker_compose ps
