#!/bin/bash
set -euo pipefail

# 服务器首次部署初始化脚本
# 在服务器上运行一次即可

DEPLOY_DIR="/opt/scare"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()  { echo -e "${GREEN}[INIT]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
err()  { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

if [ "$(id -u)" -ne 0 ]; then
    err "请使用 root 用户运行此脚本"
fi

log "===== sCare 服务器初始化 ====="

log "[1/4] 创建项目目录..."
mkdir -p "$DEPLOY_DIR/deployment/dist/c-end"
mkdir -p "$DEPLOY_DIR/deployment/dist/management-portal"

log "[2/4] 创建环境配置文件..."
ENV_FILE="$DEPLOY_DIR/deployment/.env"
if [ ! -f "$ENV_FILE" ]; then
    cat > "$ENV_FILE" << 'ENVEOF'
# sCare 生产环境配置（由 server-init.sh 生成）
# 请修改以下值为实际配置

# 数据库
DB_ROOT_PASSWORD=CHANGE_ME_root_password
DB_NAME=scare_db
DB_USER=scare_user
DB_PASSWORD=CHANGE_ME_db_password

# JWT（至少 32 字符）
JWT_SECRET=CHANGE_ME_jwt_secret_at_least_32_chars

# SMTP 邮件（可选）
SMTP_HOST=
SMTP_PORT=465
SMTP_USER=
SMTP_PASSWORD=

# 高德地图（可选）
AMAP_KEY=
ENVEOF
    warn "已生成 $ENV_FILE，请编辑填入实际密码："
    warn "  nano $ENV_FILE"
else
    log "  $ENV_FILE 已存在，跳过"
fi

log "[3/4] 配置防火墙..."
if command -v ufw &> /dev/null; then
    ufw allow 22/tcp    2>/dev/null || true
    ufw allow 80/tcp    2>/dev/null || true
    ufw allow 443/tcp   2>/dev/null || true
    log "  ufw 规则已添加"
elif command -v iptables &> /dev/null; then
    warn "  检测到 iptables，请手动配置防火墙开放 80/443 端口"
else
    warn "  未检测到防火墙工具"
fi

log "[4/4] 验证 Docker 环境..."
docker --version || err "Docker 未安装"
docker-compose --version || err "docker-compose 未安装"
log "  Docker 环境正常"

echo ""
log "===== 初始化完成 ====="
echo ""
echo "后续步骤："
echo "  1. 编辑环境变量:  nano $DEPLOY_DIR/deployment/.env"
echo "  2. 在 GitHub 仓库设置 Secrets（Settings → Secrets and variables → Actions）："
echo "     - SERVER_HOST:    你的服务器 IP"
echo "     - SERVER_USER:    SSH 用户名（建议 root）"
echo "     - SERVER_SSH_KEY: SSH 私钥"
echo "     - SERVER_PORT:    SSH 端口（默认 22 可不设）"
echo "     - AMAP_KEY:       高德地图 Key（前端构建用）"
echo "  3. 推送代码到 master 分支即可触发自动部署"
