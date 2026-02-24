#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "=== sCare 部署脚本 ==="
echo ""

if ! command -v docker &> /dev/null; then
    echo "错误: 未安装 docker"
    exit 1
fi

if ! docker compose version &> /dev/null; then
    echo "错误: 未安装 docker compose"
    exit 1
fi

if [ ! -f "$SCRIPT_DIR/.env" ]; then
    echo "错误: 未找到 $SCRIPT_DIR/.env 配置文件"
    echo "请复制 .env.example 为 .env 并填入实际配置"
    exit 1
fi

echo "[1/4] 构建 C 端前端..."
cd "$PROJECT_ROOT/frontend/c-end"
npm install --production=false
npm run build
mkdir -p "$SCRIPT_DIR/dist/c-end"
cp -r dist/* "$SCRIPT_DIR/dist/c-end/"
echo "  C 端构建完成"

echo "[2/4] 构建 B 端管理门户..."
cd "$PROJECT_ROOT/frontend/management-portal"
npm install --production=false
npm run build
mkdir -p "$SCRIPT_DIR/dist/management-portal"
cp -r dist/* "$SCRIPT_DIR/dist/management-portal/"
echo "  B 端构建完成"

echo "[3/4] 启动 Docker 服务..."
cd "$SCRIPT_DIR"
docker compose -f docker-compose.prod.yml up -d --build

echo "[4/4] 等待服务就绪..."
echo -n "  等待后端启动"
for i in $(seq 1 30); do
    if curl -sf http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        echo " ✅"
        break
    fi
    echo -n "."
    sleep 2
done

echo ""
echo "=== 部署完成 ==="
echo "  C 端:     http://localhost"
echo "  管理门户: http://localhost/manage/"
echo "  API:      http://localhost/api/v1/"
echo "  Swagger:  http://localhost/swagger/index.html"
echo ""
echo "测试账号:"
echo "  Admin:           13800000001 / Test@123"
echo "  Station Manager: 13800000002 / Test@123"
echo "  Staff:           13800000004 / Test@123"
