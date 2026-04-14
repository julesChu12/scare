#!/bin/bash
# =============================================================================
# 本地 CI 预编译测试脚本
# 模拟 GitHub Actions CI 构建流程：前端构建 + 后端 Docker 镜像构建
# 用法：./scripts/local-ci-test.sh [--full]
#   --full    完整测试：构建 + 启动 Docker Compose（默认只构建不启动）
#   --clean   清理构建产物
# =============================================================================

set -e

# ---------------------------------------------------------------------------
# 解析参数
# ---------------------------------------------------------------------------
MODE="build"  # build | full | clean

while [ $# -gt 0 ]; do
  case "$1" in
    --full)  MODE="full";  shift ;;
    --clean) MODE="clean"; shift ;;
    *)       echo "用法: $0 [--full|--clean]"; exit 1 ;;
  esac
done

if [ "$MODE" = "clean" ]; then
  echo ""
  echo -e "\033[0;34m[INFO]\033[0m  清理构建产物..."
  ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd | xargs dirname)"
  DEPLOY_DIR="$ROOT_DIR/deployment"
  rm -rf "$DEPLOY_DIR/dist/c-end"
  rm -rf "$DEPLOY_DIR/dist/management-portal"
  docker rmi scare-backend:local 2>/dev/null || true
  echo -e "\033[0;32m[ OK ]\033[0m  清理完成"
  exit 0
fi

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info()  { echo -e "${BLUE}[INFO]${NC}  $*"; }
log_ok()    { echo -e "${GREEN}[ OK ]${NC}  $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
log_fail()  { echo -e "${RED}[FAIL]${NC}  $*" >&2; }

# 目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
DEPLOY_DIR="$ROOT_DIR/deployment"
DIST_DIR="$DEPLOY_DIR/dist"

# Node.js 版本
NODE_VERSION="${NODE_VERSION:-20}"

# ---------------------------------------------------------------------------
# 辅助函数
# ---------------------------------------------------------------------------

check_tool() {
  if ! command -v "$1" &>/dev/null; then
    log_fail "缺少工具: $1，请先安装"
    exit 1
  fi
}

# ---------------------------------------------------------------------------
# Step 0: 环境检查
# ---------------------------------------------------------------------------
echo ""
log_info "========== Step 0: 环境检查 =========="

check_tool node
check_tool npm
check_tool docker

NODE_VER=$(node -v)
NPM_VER=$(npm -v)
DOCKER_VER=$(docker --version)

log_ok "Node.js $NODE_VER"
log_ok "npm $NPM_VER"
log_ok "$DOCKER_VER"
log_ok "环境检查通过"

# ---------------------------------------------------------------------------
# Step 1: 构建 C 端前端
# ---------------------------------------------------------------------------
echo ""
log_info "========== Step 1: 构建 C 端前端 (PWA) =========="

C_END_DIR="$ROOT_DIR/frontend/c-end"
C_END_DIST="$C_END_DIR/dist"

cd "$C_END_DIR"

# 安装依赖
log_info "安装 C 端依赖 (npm ci)..."
npm ci --prefer-offline --no-audit --progress=false

# 类型检查
log_info "运行 TypeScript 类型检查 (vue-tsc -b)..."
if ! npx vue-tsc -b --noEmit; then
  log_fail "C 端 TypeScript 类型检查失败"
  exit 1
fi

# 构建
log_info "执行 npm run build..."
VITE_API_BASE_URL=/api npm run build

if [ ! -d "$C_END_DIST" ]; then
  log_fail "C 端构建产物目录不存在: $C_END_DIST"
  exit 1
fi

C_END_SIZE=$(du -sh "$C_END_DIST" | cut -f1)
log_ok "C 端构建完成，产物大小: $C_END_SIZE"

# ---------------------------------------------------------------------------
# Step 2: 构建 B 端管理门户
# ---------------------------------------------------------------------------
echo ""
log_info "========== Step 2: 构建 B 端管理门户 =========="

MP_DIR="$ROOT_DIR/frontend/management-portal"
MP_DIST="$MP_DIR/dist"

cd "$MP_DIR"

# 安装依赖
log_info "安装管理门户依赖 (npm ci)..."
npm ci --prefer-offline --no-audit --progress=false

# 类型检查
log_info "运行 TypeScript 类型检查 (vue-tsc -b)..."
if ! npx vue-tsc -b --noEmit; then
  log_fail "B 端 TypeScript 类型检查失败"
  exit 1
fi

# 构建
log_info "执行 npm run build..."
VITE_API_BASE_URL=/api VITE_AMAP_KEY="${VITE_AMAP_KEY:-test_key}" npm run build

if [ ! -d "$MP_DIST" ]; then
  log_fail "管理门户构建产物目录不存在: $MP_DIST"
  exit 1
fi

MP_SIZE=$(du -sh "$MP_DIST" | cut -f1)
log_ok "管理门户构建完成，产物大小: $MP_SIZE"

# ---------------------------------------------------------------------------
# Step 3: 复制产物到 deployment/dist/
# ---------------------------------------------------------------------------
echo ""
log_info "========== Step 3: 复制产物到 deployment/dist/ =========="

mkdir -p "$DIST_DIR/c-end"
mkdir -p "$DIST_DIR/management-portal"

# 用 rsync 或 cp 复制（rsync 更高效，cp 作为 fallback）
if command -v rsync &>/dev/null; then
  rsync -a --delete "$C_END_DIST/" "$DIST_DIR/c-end/"
  rsync -a --delete "$MP_DIST/" "$DIST_DIR/management-portal/"
else
  cp -r "$C_END_DIST/"* "$DIST_DIR/c-end/"
  cp -r "$MP_DIST/"* "$DIST_DIR/management-portal/"
fi

log_ok "产物已复制到 deployment/dist/"
ls -la "$DIST_DIR/c-end/" | head -5
echo "  ..."
ls -la "$DIST_DIR/management-portal/" | head -5

# ---------------------------------------------------------------------------
# Step 4: 构建后端 Docker 镜像
# ---------------------------------------------------------------------------
echo ""
log_info "========== Step 4: 构建后端 Docker 镜像 (模拟 CI) =========="

cd "$DEPLOY_DIR"

log_info "执行 docker build -f Dockerfile.backend -t scare-backend:local ..."
if ! docker build -f Dockerfile.backend -t scare-backend:local ..; then
  log_fail "后端 Docker 镜像构建失败"
  exit 1
fi

BACKEND_IMAGE_SIZE=$(docker images scare-backend:local --format "{{.Size}}")
log_ok "后端镜像构建完成，大小: $BACKEND_IMAGE_SIZE"

# ---------------------------------------------------------------------------
# Step 5: Docker Compose 语法校验
# ---------------------------------------------------------------------------
echo ""
log_info "========== Step 5: Docker Compose 语法校验 =========="

cd "$DEPLOY_DIR"

# 仅校验 docker-compose 配置，不实际启动（避免网络拉取镜像耗时）
if docker compose -f docker-compose.prod.yml config --quiet; then
  log_ok "docker-compose.prod.yml 语法正确"
else
  log_fail "docker-compose.prod.yml 语法错误"
  exit 1
fi

# ---------------------------------------------------------------------------
# Step 6: 可选 - 启动 Docker Compose
# ---------------------------------------------------------------------------
if [ "$MODE" = "full" ]; then
  echo ""
  log_info "========== Step 6: 启动 Docker Compose (完整测试) =========="

  cd "$DEPLOY_DIR"

  # 检查 .env 文件
  if [ ! -f ".env" ]; then
    if [ -f ".env.example" ]; then
      log_warn ".env 不存在，复制 .env.example 为 .env"
      cp .env.example .env
      log_warn "请编辑 .env 填入必要的环境变量（JWT_SECRET 等）"
      log_info "启动前请先编辑 .env，然后重新运行: ./scripts/local-ci-test.sh --full"
      exit 1
    else
      log_fail ".env 文件不存在"
      exit 1
    fi
  fi

  log_info "拉取并启动服务 (docker compose up -d)..."
  docker compose -f docker-compose.prod.yml up -d

  log_info "等待服务启动..."
  sleep 10

  # 健康检查
  for i in $(seq 1 12); do
    if curl -sf http://localhost:8080/api/v1/health &>/dev/null; then
      log_ok "后端服务健康检查通过"
      break
    fi
    if [ $i -eq 12 ]; then
      log_fail "后端服务健康检查失败，请检查日志: docker compose -f docker-compose.prod.yml logs backend"
      exit 1
    fi
    echo "  等待中... ($i/12)"
    sleep 3
  done

  docker compose -f docker-compose.prod.yml ps
fi

# ---------------------------------------------------------------------------
# 完成
# ---------------------------------------------------------------------------
echo ""
log_ok "=============================================="
log_ok "  本地 CI 预编译测试全部通过！"
log_ok "=============================================="
echo ""
log_info "构建产物位于: $DIST_DIR"
log_info "Docker 镜像:   scare-backend:local"
echo ""
if [ "$MODE" = "full" ]; then
  log_info "服务已启动，请访问: http://localhost"
  log_info "停止服务: cd deployment && docker compose -f docker-compose.prod.yml down"
else
  log_info "如需完整测试（构建 + 启动服务），运行："
  log_info "  ./scripts/local-ci-test.sh --full"
fi
log_info ""
log_info "如需清理，运行："
log_info "  ./scripts/local-ci-test.sh --clean"
echo ""
