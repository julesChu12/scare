#!/bin/bash
# =====================================================
# sCare Backend API 测试脚本
# 用法: ./scripts/test_api.sh [base_url]
# 默认: http://localhost:8080
# =====================================================

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

# 调用 API 并返回响应体（用于提取 ID 等字段）
api_call() {
    local method="$1"
    local endpoint="$2"
    local token="$3"
    local data="$4"

    local auth_header=""
    if [ -n "$token" ]; then
        auth_header="-H \"Authorization: Bearer $token\""
    fi

    local data_param=""
    if [ -n "$data" ]; then
        data_param="-d '$data'"
    fi

    local cmd="curl -s -X $method \"$API_URL$endpoint\" -H \"Content-Type: application/json\" $auth_header $data_param"
    eval $cmd
}

# 登录并获取 token（B端）
login() {
    local phone="$1"
    local password="$2"
    local response=$(curl -s -X POST "$API_URL/b/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"phone\":\"$phone\",\"password\":\"$password\"}")
    echo "$response" | jq -r '.data.token // empty'
}

# 登录并获取 token（C端，密码模式）
c_login() {
    local phone="$1"
    local password="$2"
    local response=$(curl -s -X POST "$API_URL/c/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"phone\":\"$phone\",\"password\":\"$password\",\"type\":\"password\"}")
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
test_api "Admin - 任务列表" "GET" "/b/tasks" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 菜单详情(ID=1)" "GET" "/b/menus/1" "$ADMIN_TOKEN" "" ""
test_api "Admin - 新闻详情(ID=1)" "GET" "/b/news/1" "$ADMIN_TOKEN" "" ""
test_api "Admin - 统计仪表盘" "GET" "/b/statistics/dashboard" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 统计任务" "GET" "/b/statistics/tasks" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 统计请求" "GET" "/b/statistics/requests" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 统计今日" "GET" "/b/statistics/today" "$ADMIN_TOKEN" "" "ok"

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
# Admin 写操作测试（CRUD 流程）
# =====================================================
echo -e "\n${YELLOW}--- Admin 写操作测试 ---${NC}"

# --- 新闻 CRUD ---
NEWS_RESP=$(api_call "POST" "/b/news" "$ADMIN_TOKEN" '{"title":"集成测试新闻","summary":"测试摘要","content":"<p>测试正文</p>","type":"notice","status":"published","station_id":1}')
TEST_NEWS_ID=$(echo "$NEWS_RESP" | jq -r '.data.id // empty')
if [ -n "$TEST_NEWS_ID" ] && [ "$TEST_NEWS_ID" != "null" ]; then
    echo -e "${GREEN}✓${NC} Admin - 新闻创建 (ID=$TEST_NEWS_ID)"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗${NC} Admin - 新闻创建失败"
    FAILED=$((FAILED + 1))
fi
TOTAL=$((TOTAL + 1))

if [ -n "$TEST_NEWS_ID" ] && [ "$TEST_NEWS_ID" != "null" ]; then
    test_api "Admin - 新闻详情(测试)" "GET" "/b/news/$TEST_NEWS_ID" "$ADMIN_TOKEN" "" ""
    test_api "Admin - 新闻更新" "PUT" "/b/news/$TEST_NEWS_ID" "$ADMIN_TOKEN" '{"title":"集成测试新闻-已更新","summary":"更新摘要","content":"<p>更新正文</p>","type":"notice","status":"published","station_id":1}' "更新成功"
    test_api "Admin - 新闻删除" "DELETE" "/b/news/$TEST_NEWS_ID" "$ADMIN_TOKEN" "" "删除成功"
fi

# --- 轮播图 CRUD ---
BANNER_RESP=$(api_call "POST" "/b/banners" "$ADMIN_TOKEN" '{"title":"测试轮播图","image_url":"/static/test.jpg","link_type":"none","sort":99,"status":"active","station_id":1}')
TEST_BANNER_ID=$(echo "$BANNER_RESP" | jq -r '.data.id // empty')
if [ -n "$TEST_BANNER_ID" ] && [ "$TEST_BANNER_ID" != "null" ]; then
    echo -e "${GREEN}✓${NC} Admin - 轮播图创建 (ID=$TEST_BANNER_ID)"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗${NC} Admin - 轮播图创建失败"
    FAILED=$((FAILED + 1))
fi
TOTAL=$((TOTAL + 1))

if [ -n "$TEST_BANNER_ID" ] && [ "$TEST_BANNER_ID" != "null" ]; then
    test_api "Admin - 轮播图更新" "PUT" "/b/banners/$TEST_BANNER_ID" "$ADMIN_TOKEN" '{"title":"测试轮播图-已更新","image_url":"/static/test2.jpg","link_type":"none","sort":99,"status":"active","station_id":1}' "ok"
    test_api "Admin - 轮播图删除" "DELETE" "/b/banners/$TEST_BANNER_ID" "$ADMIN_TOKEN" "" "ok"
fi

# --- 站点 CRUD ---
TIMESTAMP=$(date +%s)
STATION_RESP=$(api_call "POST" "/b/stations" "$ADMIN_TOKEN" "{\"name\":\"集成测试站点_${TIMESTAMP}\",\"code\":\"TEST_${TIMESTAMP}\",\"latitude\":40.0,\"longitude\":116.4,\"status\":\"active\"}")
TEST_STATION_ID=$(echo "$STATION_RESP" | jq -r '.data.id // empty')
if [ -n "$TEST_STATION_ID" ] && [ "$TEST_STATION_ID" != "null" ]; then
    echo -e "${GREEN}✓${NC} Admin - 站点创建 (ID=$TEST_STATION_ID)"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}✗${NC} Admin - 站点创建失败"
    FAILED=$((FAILED + 1))
fi
TOTAL=$((TOTAL + 1))

if [ -n "$TEST_STATION_ID" ] && [ "$TEST_STATION_ID" != "null" ]; then
    test_api "Admin - 站点详情(测试)" "GET" "/b/stations/$TEST_STATION_ID" "$ADMIN_TOKEN" "" "ok"
    test_api "Admin - 站点更新" "PUT" "/b/stations/$TEST_STATION_ID" "$ADMIN_TOKEN" "{\"name\":\"集成测试站点_${TIMESTAMP}_已更新\",\"code\":\"TEST_${TIMESTAMP}\",\"latitude\":40.01,\"longitude\":116.41,\"status\":\"active\"}" "ok"

    # --- 围栏 CRUD（依赖站点）---
    ZONE_RESP=$(api_call "POST" "/b/zones" "$ADMIN_TOKEN" "{\"station_id\":$TEST_STATION_ID,\"name\":\"测试围栏\",\"points\":[{\"lat\":40.0,\"lng\":116.3},{\"lat\":40.1,\"lng\":116.3},{\"lat\":40.1,\"lng\":116.5},{\"lat\":40.0,\"lng\":116.5}],\"priority\":1,\"status\":\"active\"}")
    TEST_ZONE_ID=$(echo "$ZONE_RESP" | jq -r '.data.id // empty')
    if [ -n "$TEST_ZONE_ID" ] && [ "$TEST_ZONE_ID" != "null" ]; then
        echo -e "${GREEN}✓${NC} Admin - 围栏创建 (ID=$TEST_ZONE_ID)"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}✗${NC} Admin - 围栏创建失败"
        FAILED=$((FAILED + 1))
    fi
    TOTAL=$((TOTAL + 1))

    if [ -n "$TEST_ZONE_ID" ] && [ "$TEST_ZONE_ID" != "null" ]; then
        test_api "Admin - 围栏更新" "PUT" "/b/zones/$TEST_ZONE_ID" "$ADMIN_TOKEN" "{\"station_id\":$TEST_STATION_ID,\"name\":\"测试围栏-已更新\",\"points\":[{\"lat\":40.0,\"lng\":116.3},{\"lat\":40.1,\"lng\":116.3},{\"lat\":40.1,\"lng\":116.5},{\"lat\":40.0,\"lng\":116.5}],\"priority\":2,\"status\":\"active\"}" "ok"
        test_api "Admin - 围栏删除" "DELETE" "/b/zones/$TEST_ZONE_ID" "$ADMIN_TOKEN" "" "ok"
    fi

    # 删除测试站点（围栏删除后才能删站点）
    test_api "Admin - 站点删除" "DELETE" "/b/stations/$TEST_STATION_ID" "$ADMIN_TOKEN" "" "ok"
fi

# --- 用户创建（使用时间戳避免手机号冲突）---
TEST_PHONE="138${TIMESTAMP: -8}"
test_api "Admin - 创建用户(staff)" "POST" "/b/users" "$ADMIN_TOKEN" "{\"phone\":\"$TEST_PHONE\",\"password\":\"Test@123\",\"name\":\"测试员工\",\"identity_type\":\"staff\",\"station_id\":1}" "ok"

# --- 请求详情与状态更新 ---
test_api "Admin - 服务请求详情(ID=1)" "GET" "/b/requests/1" "$ADMIN_TOKEN" "" ""
test_api "Admin - 任务详情(ID=1)" "GET" "/b/tasks/1" "$ADMIN_TOKEN" "" ""

# --- 统计接口补全 ---
test_api "Admin - 统计概览" "GET" "/b/statistics/overview" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 统计服务类型" "GET" "/b/statistics/service-types" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 统计趋势" "GET" "/b/statistics/trend" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 统计效率" "GET" "/b/statistics/efficiency" "$ADMIN_TOKEN" "" "ok"
test_api "Admin - 统计员工排名" "GET" "/b/statistics/staff-ranking" "$ADMIN_TOKEN" "" "ok"

# --- 通知标记已读 ---
test_api "Admin - 通知标记已读(ID=1)" "POST" "/b/notifications/1/read" "$ADMIN_TOKEN" "" ""

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
# C端公开 API 测试
# =====================================================
echo -e "\n${YELLOW}--- C端公开 API 测试 ---${NC}"
test_api "C端 - 新闻列表" "GET" "/c/news" "" "" "ok"
test_api "C端 - 新闻详情(ID=1)" "GET" "/c/news/1" "" "" ""
test_api "C端 - Banner列表" "GET" "/c/banners" "" "" "ok"
test_api "C端 - 站点匹配" "GET" "/c/stations/match?lng=116.4&lat=39.9" "" "" ""
test_api "C端 - 站点匹配(POST)" "POST" "/c/stations/match" "" '{"latitude":39.9,"longitude":116.4}' ""
test_api "C端 - 逆地理编码" "GET" "/c/geocode/reverse?lng=116.4&lat=39.9" "" "" ""
test_api "C端 - 正地理编码" "POST" "/c/geocode" "" '{"address":"北京市昌平区霍营街道"}' ""

# =====================================================
# C端认证流程测试
# =====================================================
echo -e "\n${YELLOW}--- C端认证流程 ---${NC}"
C_TOKEN=$(c_login "13800000008" "Test@123")
if [ -z "$C_TOKEN" ]; then
    echo -e "${RED}✗${NC} C端用户(张大爷)登录失败"
    FAILED=$((FAILED + 1))
else
    echo -e "${GREEN}✓${NC} C端用户(张大爷)密码登录成功"
    PASSED=$((PASSED + 1))
fi
TOTAL=$((TOTAL + 1))

if [ -n "$C_TOKEN" ]; then
    test_api "C端 - 获取当前用户" "GET" "/c/auth/me" "$C_TOKEN" "" "ok"
    test_api "C端 - Token检查" "GET" "/c/auth/check" "$C_TOKEN" "" ""
    test_api "C端 Token - 访问B端拒绝" "GET" "/b/auth/me" "$C_TOKEN" "" "token end type mismatch" "403"
    test_api "C端 - 更新个人资料" "PUT" "/c/profile" "$C_TOKEN" '{"name":"张大爷","address":"北京市昌平区霍营街道华龙苑北里小区"}' "ok"
    test_api "C端 - 通知列表" "GET" "/c/notifications" "$C_TOKEN" "" "ok"

    # =====================================================
    # C端业务流程测试（创建请求 → 查看 → 取消）
    # =====================================================
    echo -e "\n${YELLOW}--- C端业务流程 ---${NC}"

    # 创建服务请求
    REQ_RESP=$(api_call "POST" "/c/requests" "$C_TOKEN" '{"service_type":"meal","contact_name":"张大爷","contact_phone":"13800000008","address":"北京市昌平区霍营街道华龙苑北里小区1号楼3单元501","lat":40.07,"lng":116.37}')
    C_REQUEST_ID=$(echo "$REQ_RESP" | jq -r '.data.id // empty')
    C_REQUEST_MSG=$(echo "$REQ_RESP" | jq -r '.msg // empty')
    if [ -n "$C_REQUEST_ID" ] && [ "$C_REQUEST_ID" != "null" ]; then
        echo -e "${GREEN}✓${NC} C端 - 创建服务请求 (ID=$C_REQUEST_ID)"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}✗${NC} C端 - 创建服务请求失败 (msg=$C_REQUEST_MSG)"
        FAILED=$((FAILED + 1))
    fi
    TOTAL=$((TOTAL + 1))

    # 查看请求列表和详情
    test_api "C端 - 服务请求列表" "GET" "/c/requests" "$C_TOKEN" "" "ok"

    if [ -n "$C_REQUEST_ID" ] && [ "$C_REQUEST_ID" != "null" ]; then
        test_api "C端 - 服务请求详情" "GET" "/c/requests/$C_REQUEST_ID" "$C_TOKEN" "" "ok"

        # 请求创建时已自动 dispatched 并生成任务，直接查询对应任务
        test_api "Admin - 查看C端请求详情" "GET" "/b/requests/$C_REQUEST_ID" "$ADMIN_TOKEN" "" ""

        # 通过任务列表按 request_id 找到对应任务（响应字段为 items，加大 page_size 确保查到）
        TASK_LIST_RESP=$(api_call "GET" "/b/tasks?page_size=100" "$ADMIN_TOKEN" "")
        C_TASK_ID=$(echo "$TASK_LIST_RESP" | jq -r --argjson rid "$C_REQUEST_ID" '[.data.items[] | select(.request_id == $rid)] | .[0].id // empty')

        if [ -n "$C_TASK_ID" ] && [ "$C_TASK_ID" != "null" ]; then
            # Staff 认领任务
            test_api "Staff - 认领任务" "POST" "/b/tasks/$C_TASK_ID/claim" "$STAFF_TOKEN" "" "ok"
            # Staff 完成任务
            test_api "Staff - 完成任务" "POST" "/b/tasks/$C_TASK_ID/complete" "$STAFF_TOKEN" '{"images":[]}' "ok"
            # C端评价
            test_api "C端 - 评价服务" "POST" "/c/requests/$C_REQUEST_ID/rate" "$C_TOKEN" '{"rating":5,"feedback":"服务很好"}' "ok"
        else
            echo -e "${YELLOW}⚠${NC} 未找到对应任务，跳过认领/完成/评价测试"
        fi

        # 再创建一个请求用于测试取消
        CANCEL_REQ_RESP=$(api_call "POST" "/c/requests" "$C_TOKEN" '{"service_type":"cleaning","contact_name":"张大爷","contact_phone":"13800000008","address":"北京市昌平区霍营街道","lat":40.07,"lng":116.37}')
        CANCEL_REQUEST_ID=$(echo "$CANCEL_REQ_RESP" | jq -r '.data.id // empty')
        if [ -n "$CANCEL_REQUEST_ID" ] && [ "$CANCEL_REQUEST_ID" != "null" ]; then
            test_api "C端 - 取消服务请求" "POST" "/c/requests/$CANCEL_REQUEST_ID/cancel" "$C_TOKEN" "" "ok"
            # B端管理员也可以取消请求
            CANCEL_REQ2_RESP=$(api_call "POST" "/c/requests" "$C_TOKEN" '{"service_type":"repair","contact_name":"张大爷","contact_phone":"13800000008","address":"北京市昌平区霍营街道","lat":40.07,"lng":116.37}')
            CANCEL_REQUEST_ID2=$(echo "$CANCEL_REQ2_RESP" | jq -r '.data.id // empty')
            if [ -n "$CANCEL_REQUEST_ID2" ] && [ "$CANCEL_REQUEST_ID2" != "null" ]; then
                test_api "Admin - 管理员取消请求" "POST" "/b/requests/$CANCEL_REQUEST_ID2/cancel" "$ADMIN_TOKEN" "" "ok"
            fi
        fi
    fi

    # B端更新请求信息（用新创建的 pending 请求测试）
    EDIT_REQ_RESP=$(api_call "POST" "/c/requests" "$C_TOKEN" '{"service_type":"company","contact_name":"张大爷","contact_phone":"13800000008","address":"北京市昌平区霍营街道","lat":40.07,"lng":116.37}')
    EDIT_REQUEST_ID=$(echo "$EDIT_REQ_RESP" | jq -r '.data.id // empty')
    if [ -n "$EDIT_REQUEST_ID" ] && [ "$EDIT_REQUEST_ID" != "null" ]; then
        test_api "Admin - 更新请求信息" "PUT" "/b/requests/$EDIT_REQUEST_ID" "$ADMIN_TOKEN" '{"description":"补充描述","urgency":"normal"}' ""
    fi

    # C端登出
    test_api "C端 - 登出" "POST" "/c/auth/logout" "$C_TOKEN" "" "logout successful"
fi

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
