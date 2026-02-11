#!/bin/bash

# =====================================================
# 数据库初始化脚本
# 用途: 创建数据库表结构并导入测试数据
# 使用方法: ./init.sh [选项]
# =====================================================

set -e  # 遇到错误立即退出

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 默认配置
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-secret}"
DB_NAME="${DB_NAME:-scare_db}"

# 脚本目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCHEMA_FILE="$SCRIPT_DIR/../schema/schema.sql"
SEED_FILE="$SCRIPT_DIR/../seeds/seed.sql"

# 打印帮助信息
show_help() {
    echo "数据库初始化脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help              显示帮助信息"
    echo "  -s, --schema-only       只创建表结构，不导入种子数据"
    echo "  -d, --seed-only         只导入种子数据(需表已存在)"
    echo "  -f, --force             强制重新创建数据库(会清空现有数据)"
    echo ""
    echo "环境变量:"
    echo "  DB_HOST                 数据库主机 (默认: localhost)"
    echo "  DB_PORT                 数据库端口 (默认: 3306)"
    echo "  DB_USER                 数据库用户 (默认: root)"
    echo "  DB_PASSWORD             数据库密码"
    echo "  DB_NAME                 数据库名称 (默认: scare_db)"
    echo ""
    echo "示例:"
    echo "  $0                      # 完整初始化(创建表+导入数据)"
    echo "  $0 -s                   # 只创建表结构"
    echo "  $0 -d                   # 只导入种子数据"
    echo "  DB_PASSWORD=123456 $0   # 指定密码"
}

# 打印信息
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查MySQL连接
check_mysql_connection() {
    log_info "检查MySQL连接..."

    if [ -z "$DB_PASSWORD" ]; then
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -e "SELECT 1" > /dev/null 2>&1
    else
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" > /dev/null 2>&1
    fi

    if [ $? -eq 0 ]; then
        log_info "MySQL连接成功"
    else
        log_error "MySQL连接失败，请检查配置"
        exit 1
    fi
}

# 创建数据库和表结构
create_schema() {
    log_info "开始创建数据库表结构..."

    if [ ! -f "$SCHEMA_FILE" ]; then
        log_error "表结构文件不存在: $SCHEMA_FILE"
        exit 1
    fi

    if [ -z "$DB_PASSWORD" ]; then
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" < "$SCHEMA_FILE"
    else
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < "$SCHEMA_FILE"
    fi

    if [ $? -eq 0 ]; then
        log_info "数据库表结构创建成功"
    else
        log_error "数据库表结构创建失败"
        exit 1
    fi
}

# 导入种子数据
import_seeds() {
    log_info "开始导入种子数据..."

    if [ ! -f "$SEED_FILE" ]; then
        log_error "种子数据文件不存在: $SEED_FILE"
        exit 1
    fi

    if [ -z "$DB_PASSWORD" ]; then
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" "$DB_NAME" < "$SEED_FILE"
    else
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$SEED_FILE"
    fi

    if [ $? -eq 0 ]; then
        log_info "种子数据导入成功"
    else
        log_error "种子数据导入失败"
        exit 1
    fi
}

# 强制重新创建数据库
force_recreate() {
    log_warn "警告: 即将删除现有数据库 $DB_NAME 并重新创建"
    read -p "确认继续? (yes/no): " confirm

    if [ "$confirm" != "yes" ]; then
        log_info "操作已取消"
        exit 0
    fi

    log_info "删除现有数据库..."
    if [ -z "$DB_PASSWORD" ]; then
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -e "DROP DATABASE IF EXISTS $DB_NAME;"
    else
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "DROP DATABASE IF EXISTS $DB_NAME;"
    fi
}

# 显示数据库信息
show_database_info() {
    log_info "数据库信息:"
    echo "  主机: $DB_HOST:$DB_PORT"
    echo "  数据库: $DB_NAME"
    echo "  用户: $DB_USER"
    echo ""
}

# 主函数
main() {
    SCHEMA_ONLY=false
    SEED_ONLY=false
    FORCE=false

    # 解析参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -s|--schema-only)
                SCHEMA_ONLY=true
                shift
                ;;
            -d|--seed-only)
                SEED_ONLY=true
                shift
                ;;
            -f|--force)
                FORCE=true
                shift
                ;;
            *)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done

    # 显示配置
    show_database_info

    # 检查连接
    check_mysql_connection

    # 强制重建
    if [ "$FORCE" = true ]; then
        force_recreate
    fi

    # 执行操作
    if [ "$SEED_ONLY" = true ]; then
        # 只导入种子数据
        import_seeds
    elif [ "$SCHEMA_ONLY" = true ]; then
        # 只创建表结构
        create_schema
    else
        # 完整初始化
        create_schema
        import_seeds
    fi

    log_info "数据库初始化完成！"
    echo ""
    log_info "测试账号信息:"
    echo "  管理员: 13800000001 / Test@123"
    echo "  霍营站站长: 13800000002 / Test@123"
    echo "  霍营站工作人员: 13800000004 / Test@123"
    echo "  老年人: 13800000008 / Test@123"
    echo ""
}

# 执行主函数
main "$@"
