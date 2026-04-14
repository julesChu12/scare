#!/bin/bash
# ============================================
# sCare 数据库迁移脚本
# 自动执行 database/migrations/ 下的增量 SQL
# ============================================

set -euo pipefail

# 配置
MYSQL_DB="${DB_NAME:-scare_db}"

# 脚本位于 /scripts/migrate.sh，迁移文件位于 /backend/database/migrations/
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# 判断是后端项目结构(scare/scripts/)还是旧结构(直接 scripts/)
if [ -d "$SCRIPT_DIR/../backend/database/migrations" ]; then
    MIGRATIONS_DIR="$SCRIPT_DIR/../backend/database/migrations"
else
    MIGRATIONS_DIR="$SCRIPT_DIR/database/migrations"
fi

# 在 Docker 容器内执行，使用容器内置的 MYSQL_ROOT_PASSWORD 环境变量
# 优先使用 MYSQL_ROOT_PASSWORD（容器内），否则使用 DB_ROOT_PASSWORD（.env）
if [ -n "${MYSQL_ROOT_PASSWORD:-}" ]; then
    _MYSQL_ROOT_PASS="$MYSQL_ROOT_PASSWORD"
elif [ -n "${DB_ROOT_PASSWORD:-}" ]; then
    _MYSQL_ROOT_PASS="$DB_ROOT_PASSWORD"
else
    echo -e "\033[0;31m[ERROR]\033[0m 未设置 MySQL root 密码，请设置 MYSQL_ROOT_PASSWORD 或 DB_ROOT_PASSWORD 环境变量"
    exit 1
fi
# 执行 MySQL 命令的函数
run_mysql() {
    if [ -f /.dockerenv ] || grep -q docker /proc/1/cgroup 2>/dev/null; then
        # 容器内
        mysql -u root -p"$_MYSQL_ROOT_PASS" "$MYSQL_DB" --default-character-set=utf8mb4 "$@"
    else
        # 容器外，通过 docker exec
        docker exec scare_mysql mysql -u root -p"$_MYSQL_ROOT_PASS" "$MYSQL_DB" --default-character-set=utf8mb4 "$@"
    fi
}

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[MIGRATE]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_skip()  { echo -e "${BLUE}[SKIP]${NC} $1"; }

log_info "===== sCare 数据库迁移 ====="
log_info "迁移目录: $MIGRATIONS_DIR"
log_info "MySQL 用户: root (密码来源: ${_MYSQL_ROOT_PASS:0:1}***)"

# 检查迁移目录
if [ ! -d "$MIGRATIONS_DIR" ]; then
    log_warn "迁移目录不存在，跳过迁移"
    exit 0
fi

# 测试连接
log_info "测试数据库连接..."
if ! run_mysql -e "SELECT 1" >/dev/null 2>&1; then
    log_error "无法连接数据库，请确认 MYSQL_ROOT_PASSWORD 环境变量正确"
    log_error "错误信息: $(run_mysql -e "SELECT 1" 2>&1 | grep -v "Enter password" | head -3)"
    exit 1
fi

# 创建迁移记录表
log_info "初始化迁移记录表..."
run_mysql -e "
CREATE TABLE IF NOT EXISTS \`_migrations\` (
    \`id\` INT AUTO_INCREMENT PRIMARY KEY,
    \`name\` VARCHAR(255) NOT NULL UNIQUE,
    \`applied_at\` DATETIME DEFAULT CURRENT_TIMESTAMP,
    \`checksum\` VARCHAR(64) DEFAULT NULL,
    INDEX \`idx_migrations_name\` (\`name\`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
" 2>/dev/null || {
    log_error "创建迁移记录表失败"
    exit 1
}

# 获取已执行的迁移
EXECUTED=$(run_mysql -N -e "SELECT name FROM \`_migrations\`;" 2>/dev/null | sort)

# 遍历迁移文件
MIGRATED=0
SKIPPED=0

for file in $(ls "$MIGRATIONS_DIR"/*.sql 2>/dev/null | sort); do
    filename=$(basename "$file")

    # 检查是否已执行
    if echo "$EXECUTED" | grep -qx "$filename"; then
        log_skip "$filename (已执行)"
        SKIPPED=$((SKIPPED + 1))
        continue
    fi

    log_info "执行: $filename"
    echo "-------------------------------------------"
    cat "$file"
    echo "-------------------------------------------"

    # 执行迁移
    if run_mysql < "$file"; then
        # 记录迁移
        run_mysql -e "INSERT INTO \`_migrations\` (\`name\`) VALUES ('$filename');" 2>/dev/null || true
        log_info "✓ $filename 执行成功"
        MIGRATED=$((MIGRATED + 1))
    else
        log_error "✗ $filename 执行失败"
        exit 1
    fi
done

echo ""
log_info "===== 迁移完成 ====="
if [ "$MIGRATED" -gt 0 ]; then
    log_info "新增迁移: $MIGRATED 个"
fi
if [ "$SKIPPED" -gt 0 ]; then
    log_info "跳过: $SKIPPED 个"
fi
exit 0
