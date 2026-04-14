#!/bin/bash
# ============================================
# sCare 服务器部署脚本
# 流程: 执行迁移 → 构建/启动服务
# ============================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="$SCRIPT_DIR/docker-compose.prod.yml"
ENV_FILE="$SCRIPT_DIR/.env"
MIGRATE_SCRIPT="$PROJECT_ROOT/scripts/migrate.sh"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()  { echo -e "${GREEN}[DEPLOY]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()  { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

log "===== sCare 服务器部署 ====="
log "时间: $(date '+%Y-%m-%d %H:%M:%S')"

if [ ! -f "$ENV_FILE" ]; then
    err "未找到 $ENV_FILE，请先创建环境配置文件（参考 .env.example）"
fi

# 加载环境变量
source "$ENV_FILE"

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

# ============================================
# 步骤 1: 执行数据库迁移
# ============================================
log "[1/4] 执行数据库迁移..."
if [ -f "$MIGRATE_SCRIPT" ]; then
    chmod +x "$MIGRATE_SCRIPT"
    cd "$PROJECT_ROOT"
    bash "$MIGRATE_SCRIPT" || warn "迁移过程中有警告"
else
    warn "迁移脚本不存在，跳过"
fi

# ============================================
# 步骤 2: 停止旧服务并修复目录权限
# ============================================
log "[2/4] 停止旧服务..."
cd "$SCRIPT_DIR"
docker-compose -f "$COMPOSE_FILE" down --remove-orphans 2>/dev/null || true

# 修复 backend/database 目录权限（Docker init 脚本可能以 root 创建）
if [ -d "$PROJECT_ROOT/backend/database" ]; then
    chmod -R 755 "$PROJECT_ROOT/backend/database" 2>/dev/null || true
fi

# ============================================
# 步骤 3: 构建并启动服务
# ============================================
log "[3/4] 构建并启动服务..."
docker-compose -f "$COMPOSE_FILE" up -d --build

# ============================================
# 步骤 4: 健康检查
# ============================================
log "[4/4] 健康检查..."
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
    docker-compose -f "$COMPOSE_FILE" logs --tail=50 backend
    exit 1
fi

log "===== 部署完成 ====="
log "  后端 API:  http://localhost:8080"
log "  Nginx:     http://localhost"
log "  管理后台:  http://localhost/manage/"
echo ""
docker-compose -f "$COMPOSE_FILE" ps
