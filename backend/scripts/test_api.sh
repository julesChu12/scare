#!/bin/bash
# =====================================================
# sCare Backend API 测试脚本
# 用法: ./scripts/test_api.sh [base_url]
# 默认: http://localhost:8080
# =====================================================

set -e

BASE_URL="${1:-http://localhost:8080}"
API_URL="$BASE_URL/api/v1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试计数
TOTAL=0
PASSED=0
FAILED=0

# 测试函数
test_api() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local token="$4"
    local data="$5"
    local expected_msg="$6"
    local expected_code="${7:-200}"

    TOTAL=$((TOTAL + 1))

    local auth_header=""
    if [ -n "$token" ]; then
        auth_header="-H \"Authorization: Bearer $token\""
    fi

    local data_param=""
    if [ -n "$data" ]; then
        data_param="-d '$data'"
    fi

    local cmd="curl -s -X $method \"$API_URL$endpoint\" -H \"Content-Type: application/json\" $auth_header $data_param -w '\n%{http_code}'"
    local raw_response=$(eval $cmd)
    local http_code=$(echo "$raw_response" | tail -n 1)
    local response=$(echo "$raw_response" | sed '$d')
    local msg=$(echo "$response" | jq -r '.msg // empty' 2>/dev/null)

    if [ "$http_code" = "$expected_code" ] && { [ "$msg" = "$expected_msg" ] || [ -z "$expected_msg" ]; }; then
        echo -e "${GREEN}✓${NC} $name"
        PASSED=$((PASSED + 1))
        return 0
    else
        echo -e "${RED}✗${NC} $name (expected code/msg: $expected_code/$expected_msg, got: $http_code/$msg)"
        echo "响应体: $response"
        FAILED=$((FAILED + 1))
        return 1
    fi
}

assert_not_empty() {
    local name="$1"
    local value="$2"

    TOTAL=$((TOTAL + 1))
    if [ -n "$value" ] && [ "$value" != "null" ]; then
        echo -e "${GREEN}✓${NC} $name"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}✗${NC} $name (value is empty)"
        FAILED=$((FAILED + 1))
    fi
}

# 登录并获取 token
login() {
    local phone="$1"
    local password="$2"
    local response=$(curl -s -X POST "$API_URL/b/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"phone\":\"$phone\",\"password\":\"$password\"}")
    echo "$response" | jq -r '.data.token // empty'
}

echo "=========================================="
echo "sCare Backend API 测试"
echo "Base URL: $BASE_URL"
echo "=========================================="

# =====================================================
# 健康检查
# =====================================================
echo -e "\n${YELLOW}--- 健康检查 ---${NC}"
HEALTH=$(curl -s "$API_URL/health" | jq -r '.msg // empty')
if [ "$HEALTH" = "ok" ]; then
    echo -e "${GREEN}✓${NC} 服务健康检查"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗${NC} 服务健康检查失败"
    FAILED=$((FAILED + 1))
    echo "服务可能未启动，请先运行: air"
    exit 1
fi
TOTAL=$((TOTAL + 1))

# 鉴权与端类型隔离基础校验
test_api "未认证 - B端 me 拒绝" "GET" "/b/auth/me" "" "" "missing authorization" "401"

# =====================================================
# Admin 角色测试
# =====================================================
echo -e "\n${YELLOW}--- Admin 角色测试 ---${NC}"
ADMIN_TOKEN=$(login "13800000001" "Test@123")
if [ -z "$ADMIN_TOKEN" ]; then
    echo -e "${RED}✗${NC} Admin 登录失败"
    exit 1
fi
echo -e "${GREEN}✓${NC} Admin 登录成功"
PASSED=$((PASSED + 1))
TOTAL=$((TOTAL + 1))

test_api "Admin - 获取当前用户" "GET" "/b/auth/me" "$ADMIN_TOKEN" "" "ok"
test_api "Admin Token - 访问 C端 me 拒绝" "GET" "/c/auth/me" "$ADMIN_TOKEN" "" "token end type mismatch" "403"
test_api "Admin - 用户列表" "GET" "/b/users" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 用户详情" "GET" "/b/users/1" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 用户更新(设置身份证)" "PUT" "/b/users/1" "$ADMIN_TOKEN" "{\"id_card\":\"110101199001011234\",\"name\":\"系统管理员\"}" "ok"
test_api "Admin - 用户筛选(role=staff)" "GET" "/b/users?role=staff" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 用户筛选(station_id=1)" "GET" "/b/users?station_id=1" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 用户菜单" "GET" "/b/menus/user" "$ADMIN_TOKEN" "" "success"
test_api "Admin - 菜单树" "GET" "/b/menus" "$ADMIN_TOKEN" "" "success"
test_api "Admin - 站点列表" "GET" "/b/stations" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 站点详情" "GET" "/b/stations/1" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 围栏列表" "GET" "/b/zones" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 任务池" "GET" "/b/tasks/pool" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 任务池(所有站点)" "GET" "/b/tasks/pool?station_id=0" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 任务池(指定站点)" "GET" "/b/tasks/pool?station_id=1" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 我的任务" "GET" "/b/tasks/my" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 服务请求列表" "GET" "/b/requests" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - Banner列表" "GET" "/b/banners" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 新闻列表" "GET" "/b/news" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 权限树" "GET" "/b/permissions/tree" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 角色权限" "GET" "/b/roles/admin/permissions" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 通知列表" "GET" "/b/notifications" "$ADMIN_TOKEN" "" "ok"

# 用户更新新增字段与校验分支回归
USER_DETAIL=$(curl -s -X GET "$API_URL/b/users/1" -H "Authorization: Bearer $ADMIN_TOKEN")
ID_CARD_HASH=$(echo "$USER_DETAIL" | jq -r '.data.id_card_hash // empty')
ID_CARD_TOKEN=$(echo "$USER_DETAIL" | jq -r '.data.id_card_token // empty')
ID_CARD_MASKED=$(echo "$USER_DETAIL" | jq -r '.data.id_card_masked // empty')

assert_not_empty "Admin - 用户详情返回 id_card_hash" "$ID_CARD_HASH"
assert_not_empty "Admin - 用户详情返回 id_card_token" "$ID_CARD_TOKEN"
assert_not_empty "Admin - 用户详情返回 id_card_masked" "$ID_CARD_MASKED"

test_api "Admin - 用户更新(身份证+hash互斥校验)" "PUT" "/b/users/1" "$ADMIN_TOKEN" "{\"id_card\":\"110101199001011234\",\"id_card_hash\":\"$ID_CARD_HASH\"}" "id_card_hash should not be sent when id_card is provided" "400"
test_api "Admin - 用户更新(错误id_card_token)" "PUT" "/b/users/1" "$ADMIN_TOKEN" "{\"id_card_token\":\"invalid.token\"}" "invalid id_card_token" "400"
test_api "Admin - 用户更新(id_card_hash保持原值)" "PUT" "/b/users/1" "$ADMIN_TOKEN" "{\"id_card_hash\":\"$ID_CARD_HASH\",\"name\":\"系统管理员\"}" "ok" "200"
test_api "Admin - 用户更新(id_card_token保持原值)" "PUT" "/b/users/1" "$ADMIN_TOKEN" "{\"id_card_token\":\"$ID_CARD_TOKEN\",\"name\":\"系统管理员\"}" "ok" "200"

# =====================================================
# Station Manager 角色测试
# =====================================================
echo -e "\n${YELLOW}--- Station Manager 角色测试 ---${NC}"
SM_TOKEN=$(login "13800000002" "Test@123")
if [ -z "$SM_TOKEN" ]; then
    echo -e "${RED}✗${NC} Station Manager 登录失败"
else
    echo -e "${GREEN}✓${NC} Station Manager 登录成功"
    PASSED=$((PASSED + 1))
fi
TOTAL=$((TOTAL + 1))

test_api "SM - 获取当前用户" "GET" "/b/auth/me" "$SM_TOKEN" "" "ok"
test_api "SM - 用户菜单" "GET" "/b/menus/user" "$SM_TOKEN" "" "success"
test_api "SM - 站点列表" "GET" "/b/stations" "$SM_TOKEN" "" "ok"
test_api "SM - 任务池" "GET" "/b/tasks/pool" "$SM_TOKEN" "" "ok"
test_api "SM - 服务请求列表" "GET" "/b/requests" "$SM_TOKEN" "" "ok"
test_api "SM - 用户列表" "GET" "/b/users" "$SM_TOKEN" "" "ok"

# =====================================================
# Staff 角色测试
# =====================================================
echo -e "\n${YELLOW}--- Staff 角色测试 ---${NC}"
STAFF_TOKEN=$(login "13800000004" "Test@123")
if [ -z "$STAFF_TOKEN" ]; then
    echo -e "${RED}✗${NC} Staff 登录失败"
else
    echo -e "${GREEN}✓${NC} Staff 登录成功"
    PASSED=$((PASSED + 1))
fi
TOTAL=$((TOTAL + 1))

test_api "Staff - 获取当前用户" "GET" "/b/auth/me" "$STAFF_TOKEN" "" "ok"
test_api "Staff - 用户菜单" "GET" "/b/menus/user" "$STAFF_TOKEN" "" "success"
test_api "Staff - 任务池" "GET" "/b/tasks/pool" "$STAFF_TOKEN" "" "ok"
test_api "Staff - 我的任务" "GET" "/b/tasks/my" "$STAFF_TOKEN" "" "ok"
test_api "Staff - 角色权限(无权限)" "GET" "/b/roles/admin/permissions" "$STAFF_TOKEN" "" "forbidden" "403"

# =====================================================
# C端 API 测试
# =====================================================
echo -e "\n${YELLOW}--- C端公开 API 测试 ---${NC}"
test_api "C端 - 新闻列表" "GET" "/c/news" "" "" "ok"
test_api "C端 - Banner列表" "GET" "/c/banners" "" "" "ok"

# =====================================================
# 测试结果汇总
# =====================================================
echo -e "\n=========================================="
echo "测试结果汇总"
echo "=========================================="
echo -e "总计: $TOTAL"
echo -e "${GREEN}通过: $PASSED${NC}"
echo -e "${RED}失败: $FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}所有测试通过!${NC}"
    exit 0
else
    echo -e "\n${RED}有 $FAILED 个测试失败${NC}"
    exit 1
fi
